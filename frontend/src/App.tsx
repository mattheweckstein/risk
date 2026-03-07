import { useState, useCallback, useRef } from 'react';
import { GameState } from './types/game';
import { useGameApi } from './hooks/useGameApi';
import SetupScreen from './components/SetupScreen';
import Map from './components/Map';
import GameControls from './components/GameControls';
import HistoryPanel from './components/HistoryPanel';
import ToastContainer, { ToastMessage } from './components/Toast';

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export default function App() {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [selectedTerritory, setSelectedTerritory] = useState<string | null>(null);
  const [attackTarget, setAttackTarget] = useState<string | null>(null);
  const [fortifySource, setFortifySource] = useState<string | null>(null);
  const [fortifyTarget, setFortifyTarget] = useState<string | null>(null);
  const [aiThinking, setAiThinking] = useState(false);
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const toastIdRef = useRef(0);
  const api = useGameApi();

  const addToast = useCallback((text: string, type: ToastMessage['type'] = 'info') => {
    const id = ++toastIdRef.current;
    setToasts((prev) => [...prev.slice(-4), { id, text, type }]);
  }, []);

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const clearSelection = useCallback(() => {
    setSelectedTerritory(null);
    setAttackTarget(null);
    setFortifySource(null);
    setFortifyTarget(null);
  }, []);

  // Run AI turns until it's the human player's turn
  const runAiTurns = useCallback(
    async (state: GameState) => {
      let current = state;
      const humanId = current.players.find((p) => !p.isAI)?.id;

      while (
        current.phase !== 'ended' &&
        current.currentPlayer !== humanId
      ) {
        setAiThinking(true);
        await delay(600);
        try {
          current = await api.aiTurn(current.id);
          setGameState(current);
        } catch (e) {
          console.error('AI turn failed:', e);
          break;
        }
      }
      setAiThinking(false);
      return current;
    },
    [api],
  );

  const handleGameStart = useCallback(
    async (game: GameState) => {
      setGameState(game);
      clearSelection();
      // If game starts with AI turn (setup), run AI turns
      const humanId = game.players.find((p) => !p.isAI)?.id;
      if (game.currentPlayer !== humanId) {
        await runAiTurns(game);
      }
    },
    [clearSelection, runAiTurns],
  );

  const handleTerritoryClick = useCallback(
    async (territoryId: string) => {
      if (!gameState || aiThinking) return;

      const territory = gameState.territories[territoryId];
      if (!territory) return;

      const currentPlayer = gameState.players.find((p) => p.id === gameState.currentPlayer);
      if (!currentPlayer || currentPlayer.isAI) return;

      const humanId = currentPlayer.id;

      // === SETUP PHASE ===
      if (gameState.phase === 'setup') {
        if (territory.owner !== humanId) {
          addToast('You can only place troops on your own territories', 'warning');
          return;
        }
        try {
          let updated = await api.placeTroops(gameState.id, territoryId, 1);
          setGameState(updated);

          // After human places, run AI setup turns
          if (updated.currentPlayer !== humanId && updated.phase !== 'ended') {
            updated = await runAiTurns(updated);
          }
        } catch (e) {
          addToast(e instanceof Error ? e.message : 'Failed to place troops', 'warning');
        }
        return;
      }

      // === PLACE PHASE ===
      if (gameState.phase === 'place') {
        if (territory.owner !== humanId) {
          addToast('You can only deploy to your own territories', 'warning');
          return;
        }
        if (gameState.troopsToDeploy <= 0) return;
        try {
          const updated = await api.placeTroops(gameState.id, territoryId, 1);
          setGameState(updated);
          if (updated.troopsToDeploy === 0) {
            addToast('All troops deployed! Attack phase begins.', 'success');
          }
        } catch (e) {
          addToast(e instanceof Error ? e.message : 'Failed to place troops', 'warning');
        }
        return;
      }

      // === ATTACK PHASE ===
      if (gameState.phase === 'attack') {
        // If no territory selected, select an owned territory with 2+ troops
        if (!selectedTerritory) {
          if (territory.owner !== humanId) {
            addToast('Select one of your territories to attack from', 'warning');
            return;
          }
          if (territory.troops < 2) {
            addToast('Need at least 2 troops to attack from this territory', 'warning');
            return;
          }
          setSelectedTerritory(territoryId);
          setAttackTarget(null);
          return;
        }

        // If clicking the same selected territory, deselect
        if (selectedTerritory === territoryId) {
          clearSelection();
          return;
        }

        // If clicking another owned territory, switch selection
        if (territory.owner === humanId) {
          if (territory.troops < 2) {
            addToast('Need at least 2 troops to attack from this territory', 'warning');
            return;
          }
          setSelectedTerritory(territoryId);
          setAttackTarget(null);
          return;
        }

        // Clicking an enemy territory - check if adjacent
        const sourceTerr = gameState.territories[selectedTerritory];
        if (!sourceTerr?.neighbors.includes(territoryId)) {
          addToast('Target must be adjacent to your territory', 'warning');
          return;
        }

        // Set as attack target (dice selection shown in controls panel)
        setAttackTarget(territoryId);
        return;
      }

      // === FORTIFY PHASE ===
      if (gameState.phase === 'fortify') {
        if (!fortifySource) {
          // Select source
          if (territory.owner !== humanId) {
            addToast('Select one of your territories to move troops from', 'warning');
            return;
          }
          if (territory.troops < 2) {
            addToast('Need at least 2 troops to fortify from this territory', 'warning');
            return;
          }
          setFortifySource(territoryId);
          setFortifyTarget(null);
          return;
        }

        // Clicking same source deselects
        if (fortifySource === territoryId) {
          clearSelection();
          return;
        }

        // Clicking another owned territory as target
        if (territory.owner === humanId && territoryId !== fortifySource) {
          setFortifyTarget(territoryId);
          return;
        }

        // Clicking enemy territory in fortify - ignore
        if (territory.owner !== humanId) {
          addToast('You can only fortify to your own territories', 'warning');
          return;
        }
      }
    },
    [gameState, aiThinking, selectedTerritory, attackTarget, fortifySource, api, addToast, clearSelection, runAiTurns],
  );

  const handleAttack = useCallback(
    async (dice: number) => {
      if (!gameState || !selectedTerritory || !attackTarget) return;
      try {
        const updated = await api.attack(gameState.id, selectedTerritory, attackTarget, dice);
        setGameState(updated);

        if (updated.lastAttackResult?.conquered) {
          addToast(
            `Conquered ${updated.territories[attackTarget]?.name}!`,
            'success'
          );
          // Clear attack target after conquering
          setAttackTarget(null);
          // If the source territory now has < 2 troops, deselect it too
          if ((updated.territories[selectedTerritory]?.troops || 0) < 2) {
            setSelectedTerritory(null);
          }
        }

        if (updated.phase === 'ended') {
          addToast(
            `${updated.players.find((p) => p.id === updated.winner)?.name} wins the game!`,
            'success'
          );
          clearSelection();
        }
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Attack failed', 'warning');
      }
    },
    [gameState, selectedTerritory, attackTarget, api, addToast, clearSelection],
  );

  const handleEndPhase = useCallback(async () => {
    if (!gameState) return;
    try {
      let updated = await api.endPhase(gameState.id);
      setGameState(updated);
      clearSelection();

      // Check if next player is AI
      const humanId = updated.players.find((p) => !p.isAI)?.id;
      if (updated.currentPlayer !== humanId && updated.phase !== 'ended') {
        updated = await runAiTurns(updated);
      }
    } catch (e) {
      addToast(e instanceof Error ? e.message : 'Failed to end phase', 'warning');
    }
  }, [gameState, api, clearSelection, runAiTurns, addToast]);

  const handleFortify = useCallback(
    async (troops: number) => {
      if (!gameState || !fortifySource || !fortifyTarget) return;
      try {
        let updated = await api.fortify(gameState.id, fortifySource, fortifyTarget, troops);
        setGameState(updated);
        clearSelection();
        addToast(`Moved ${troops} troops`, 'info');

        // After fortify, end the phase
        updated = await api.endPhase(updated.id);
        setGameState(updated);

        // Check if next player is AI
        const humanId = updated.players.find((p) => !p.isAI)?.id;
        if (updated.currentPlayer !== humanId && updated.phase !== 'ended') {
          updated = await runAiTurns(updated);
        }
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Fortify failed', 'warning');
      }
    },
    [gameState, fortifySource, fortifyTarget, api, clearSelection, runAiTurns, addToast],
  );

  const handleTradeCards = useCallback(
    async (indices: number[]) => {
      if (!gameState) return;
      try {
        const updated = await api.tradeCards(gameState.id, indices);
        setGameState(updated);
        addToast('Cards traded for bonus troops!', 'success');
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Trade failed', 'warning');
      }
    },
    [gameState, api, addToast],
  );

  // Compute valid targets for visual feedback
  const validTargets: string[] = [];
  if (gameState) {
    const humanId = gameState.players.find((p) => !p.isAI)?.id;

    if (gameState.phase === 'attack' && selectedTerritory && !attackTarget) {
      const source = gameState.territories[selectedTerritory];
      if (source) {
        for (const nid of source.neighbors) {
          const neighbor = gameState.territories[nid];
          if (neighbor && neighbor.owner !== humanId) {
            validTargets.push(nid);
          }
        }
      }
    }

    if (gameState.phase === 'fortify' && fortifySource && !fortifyTarget) {
      // Show all owned territories as potential targets (actual connection check done by backend)
      for (const [tid, t] of Object.entries(gameState.territories)) {
        if (t.owner === humanId && tid !== fortifySource) {
          validTargets.push(tid);
        }
      }
    }
  }

  // === RENDER ===

  if (!gameState) {
    return <SetupScreen onGameStart={handleGameStart} />;
  }

  return (
    <div className="h-screen flex flex-col overflow-hidden" style={{ background: '#1a1a2e' }}>
      <ToastContainer toasts={toasts} onRemove={removeToast} />

      {/* AI thinking overlay */}
      {aiThinking && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div
            className="px-8 py-4 rounded-xl shadow-2xl border"
            style={{ background: '#16213e', borderColor: 'rgba(233, 69, 96, 0.3)' }}
          >
            <div className="flex items-center gap-3">
              <div className="w-5 h-5 border-2 border-red-500 border-t-transparent rounded-full animate-spin" />
              <span className="text-white font-semibold">AI is thinking...</span>
            </div>
          </div>
        </div>
      )}

      {/* Game over overlay */}
      {gameState.phase === 'ended' && (
        <div className="fixed inset-0 z-40 flex items-center justify-center bg-black/60">
          <div
            className="text-center px-12 py-8 rounded-2xl shadow-2xl border"
            style={{
              background: 'linear-gradient(135deg, #16213e, #0f3460)',
              borderColor: 'rgba(255, 215, 0, 0.4)',
            }}
          >
            <div className="text-5xl font-bold text-yellow-400 mb-3">VICTORY</div>
            <div className="text-xl text-gray-300 mb-6">
              {gameState.players.find((p) => p.id === gameState.winner)?.name} has conquered the world!
            </div>
            <button
              onClick={() => {
                setGameState(null);
                clearSelection();
              }}
              className="px-8 py-3 rounded-lg font-bold text-lg uppercase tracking-wider transition-all hover:scale-105"
              style={{
                background: 'linear-gradient(135deg, #e94560, #c23152)',
                color: 'white',
              }}
            >
              New Game
            </button>
          </div>
        </div>
      )}

      {/* Main layout */}
      <div className="flex flex-1 min-h-0">
        {/* Map - 70% */}
        <div className="flex-[7] min-w-0 p-2">
          <Map
            gameState={gameState}
            selectedTerritory={selectedTerritory || fortifySource}
            validTargets={validTargets}
            onTerritoryClick={handleTerritoryClick}
          />
        </div>

        {/* Side panel - 30% */}
        <div className="flex-[3] min-w-[280px] max-w-[380px] flex flex-col border-l border-gray-700/30">
          {/* Controls - top portion */}
          <div className="flex-[6] min-h-0 overflow-y-auto">
            <GameControls
              gameState={gameState}
              onEndPhase={handleEndPhase}
              onTradeCards={handleTradeCards}
              selectedTerritory={selectedTerritory}
              attackTarget={attackTarget}
              onAttack={handleAttack}
              onFortify={handleFortify}
              fortifySource={fortifySource}
              fortifyTarget={fortifyTarget}
            />
          </div>

          {/* History - bottom portion */}
          <div className="flex-[4] min-h-0 border-t border-gray-700/30">
            <HistoryPanel log={gameState.log} />
          </div>
        </div>
      </div>
    </div>
  );
}
