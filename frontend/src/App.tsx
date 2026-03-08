import { useState, useCallback, useRef, useEffect } from 'react';
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

const GAME_ID_KEY = 'risk-game-id';

export default function App() {
  const [gameState, setGameState] = useState<GameState | null>(null);
  const [selectedTerritory, setSelectedTerritory] = useState<string | null>(null);
  const [attackTarget, setAttackTarget] = useState<string | null>(null);
  const [fortifySource, setFortifySource] = useState<string | null>(null);
  const [fortifyTarget, setFortifyTarget] = useState<string | null>(null);
  const [aiThinking, setAiThinking] = useState(false);
  const [attackPending, setAttackPending] = useState(false);
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const [instructionText, setInstructionText] = useState('');
  const [loading, setLoading] = useState(true);
  const [lastPlacedTerritory, setLastPlacedTerritory] = useState<string | null>(null);
  const toastIdRef = useRef(0);
  const api = useGameApi();

  const addToast = useCallback((text: string, type: ToastMessage['type'] = 'info') => {
    const id = ++toastIdRef.current;
    setToasts((prev) => [...prev.slice(-4), { id, text, type }]);
  }, []);

  const removeToast = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  // Restore game from localStorage on mount
  useEffect(() => {
    const savedId = localStorage.getItem(GAME_ID_KEY);
    if (savedId) {
      api.getGame(savedId)
        .then((game) => {
          setGameState(game);
          setLoading(false);
        })
        .catch(() => {
          localStorage.removeItem(GAME_ID_KEY);
          setLoading(false);
        });
    } else {
      setLoading(false);
    }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Persist game ID to localStorage
  useEffect(() => {
    if (gameState?.id) {
      localStorage.setItem(GAME_ID_KEY, gameState.id);
    }
  }, [gameState?.id]);

  const clearSelection = useCallback(() => {
    setSelectedTerritory(null);
    setAttackTarget(null);
    setFortifySource(null);
    setFortifyTarget(null);
  }, []);

  // Update instruction text based on game state
  useEffect(() => {
    if (!gameState || aiThinking) {
      setInstructionText('');
      return;
    }
    const currentPlayer = gameState.players.find((p) => p.id === gameState.currentPlayer);
    if (!currentPlayer || currentPlayer.isAI) {
      setInstructionText('');
      return;
    }

    if (gameState.phase === 'setup') {
      setInstructionText('Click one of your territories to place 1 troop');
    } else if (gameState.phase === 'place') {
      setInstructionText(`Click your territories to deploy troops (${gameState.troopsToDeploy} remaining)`);
    } else if (gameState.phase === 'attack') {
      if (gameState.pendingConquest) {
        setInstructionText(`Conquered! Move additional troops in (up to ${gameState.pendingConquest.maxTroops} more)`);
      } else if (!selectedTerritory) {
        setInstructionText('Select a territory with 2+ troops to attack from, or End Attack');
      } else if (!attackTarget) {
        setInstructionText('Now click an adjacent enemy territory to attack');
      } else {
        setInstructionText('Choose number of dice or click attack target again to auto-attack');
      }
    } else if (gameState.phase === 'fortify') {
      if (!fortifySource) {
        setInstructionText('Select a territory to move troops from, or Skip');
      } else if (!fortifyTarget) {
        setInstructionText('Select a connected friendly territory to move troops to');
      } else {
        setInstructionText('Choose how many troops to move');
      }
    } else if (gameState.phase === 'ended') {
      setInstructionText('');
    }
  }, [gameState, aiThinking, selectedTerritory, attackTarget, fortifySource, fortifyTarget]);

  // Run AI turns until it's the human player's turn
  const runAiTurns = useCallback(
    async (state: GameState) => {
      let current = state;
      const humanId = current.players.find((p) => !p.isAI)?.id;

      while (current.phase !== 'ended' && current.currentPlayer !== humanId) {
        setAiThinking(true);
        await delay(500);
        try {
          current = await api.aiTurn(current.id);
          setGameState(current);
        } catch (e) {
          console.error('AI turn failed:', e);
          addToast('AI turn failed - check console', 'warning');
          break;
        }
      }
      setAiThinking(false);
      return current;
    },
    [api, addToast],
  );

  const handleGameStart = useCallback(
    async (game: GameState) => {
      setGameState(game);
      clearSelection();
      const humanId = game.players.find((p) => !p.isAI)?.id;
      if (game.currentPlayer !== humanId) {
        await runAiTurns(game);
      }
    },
    [clearSelection, runAiTurns],
  );

  // Deploy all remaining troops to a territory
  const handleDeployAll = useCallback(
    async (territoryId: string) => {
      if (!gameState || gameState.troopsToDeploy <= 0) return;
      try {
        const updated = await api.placeTroops(gameState.id, territoryId, gameState.troopsToDeploy);
        setGameState(updated);
        addToast('All troops deployed!', 'success');
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Failed to deploy', 'warning');
      }
    },
    [gameState, api, addToast],
  );

  // Execute attack with max dice by default
  const doAttack = useCallback(
    async (from: string, to: string, dice?: number) => {
      if (!gameState) return;
      // Block if pending conquest exists
      if (gameState.pendingConquest) {
        addToast('Move troops into conquered territory first', 'warning');
        return;
      }
      const fromTerr = gameState.territories[from];
      const maxDice = Math.min(3, (fromTerr?.troops || 1) - 1);
      const attackDice = dice || maxDice;
      if (attackDice < 1) return;

      setAttackPending(true);
      try {
        const updated = await api.attack(gameState.id, from, to, attackDice);
        setGameState(updated);

        if (updated.lastAttackResult?.conquered) {
          addToast(`Conquered ${updated.territories[to]?.name}!`, 'success');
          // If there's a pending conquest (can move more troops), keep UI in conquest mode
          // Otherwise clear targets
          if (!updated.pendingConquest) {
            setAttackTarget(null);
            const newSourceTroops = updated.territories[from]?.troops || 0;
            if (newSourceTroops < 2) {
              setSelectedTerritory(null);
            }
          }
        } else {
          const newSourceTroops = updated.territories[from]?.troops || 0;
          if (newSourceTroops < 2) {
            clearSelection();
          }
        }

        if (updated.phase === 'ended') {
          const winner = updated.players.find((p) => p.id === updated.winner);
          addToast(`${winner?.name} wins the game!`, 'success');
          clearSelection();
        }
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Attack failed', 'warning');
      } finally {
        setAttackPending(false);
      }
    },
    [gameState, api, addToast, clearSelection],
  );

  // Blitz: auto-attack with max dice until win or can't continue
  const handleBlitz = useCallback(
    async () => {
      if (!gameState || !selectedTerritory || !attackTarget) return;
      const from = selectedTerritory;
      const to = attackTarget;
      setAttackPending(true);
      try {
        let current = gameState;
        while (true) {
          const fromTerr = current.territories[from];
          const toTerr = current.territories[to];
          if (!fromTerr || !toTerr) break;
          const maxDice = Math.min(3, (fromTerr.troops || 1) - 1);
          if (maxDice < 1) break;
          // Stop if we already own the target (conquered)
          if (toTerr.owner === current.players.find((p) => !p.isAI)?.id) break;

          const updated = await api.attack(current.id, from, to, maxDice);
          current = updated;
          setGameState(updated);

          if (updated.lastAttackResult?.conquered) {
            addToast(`Conquered ${updated.territories[to]?.name}!`, 'success');
            if (!updated.pendingConquest) {
              setAttackTarget(null);
              if ((updated.territories[from]?.troops || 0) < 2) setSelectedTerritory(null);
            }
            break;
          }
          if (updated.phase === 'ended') {
            const winner = updated.players.find((p) => p.id === updated.winner);
            addToast(`${winner?.name} wins the game!`, 'success');
            clearSelection();
            break;
          }
        }
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Blitz failed', 'warning');
      } finally {
        setAttackPending(false);
      }
    },
    [gameState, selectedTerritory, attackTarget, api, addToast, clearSelection],
  );

  // Move additional troops after conquest
  const handleConquestMove = useCallback(
    async (troops: number) => {
      if (!gameState) return;
      try {
        const updated = await api.moveAfterConquest(gameState.id, troops);
        setGameState(updated);
        // Clear attack selection after moving troops
        setAttackTarget(null);
        if (selectedTerritory) {
          const newSourceTroops = updated.territories[selectedTerritory]?.troops || 0;
          if (newSourceTroops < 2) {
            setSelectedTerritory(null);
          }
        }
      } catch (e) {
        addToast(e instanceof Error ? e.message : 'Failed to move troops', 'warning');
      }
    },
    [gameState, api, addToast, selectedTerritory],
  );

  const handleTerritoryClick = useCallback(
    async (territoryId: string) => {
      if (!gameState || aiThinking || attackPending) return;

      // Block territory clicks during pending conquest
      if (gameState.pendingConquest) {
        addToast('Move troops into conquered territory first (use panel on right)', 'warning');
        return;
      }

      const territory = gameState.territories[territoryId];
      if (!territory) return;

      const currentPlayer = gameState.players.find((p) => p.id === gameState.currentPlayer);
      if (!currentPlayer || currentPlayer.isAI) return;

      const humanId = currentPlayer.id;

      // === SETUP PHASE ===
      if (gameState.phase === 'setup') {
        if (territory.owner !== humanId) {
          addToast('Click one of YOUR territories (highlighted)', 'warning');
          return;
        }
        try {
          let updated = await api.placeTroops(gameState.id, territoryId, 1);
          setGameState(updated);
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
          setLastPlacedTerritory(territoryId);
          if (updated.troopsToDeploy === 0) {
            setLastPlacedTerritory(null);
            addToast('All troops deployed! Now attack or end attack phase.', 'success');
          }
        } catch (e) {
          addToast(e instanceof Error ? e.message : 'Failed to place troops', 'warning');
        }
        return;
      }

      // === ATTACK PHASE ===
      if (gameState.phase === 'attack') {
        // Clicking own territory
        if (territory.owner === humanId) {
          if (territory.troops < 2) {
            addToast('Need at least 2 troops to attack from here', 'warning');
            return;
          }
          // Select/switch attacker
          if (selectedTerritory === territoryId) {
            clearSelection();
          } else {
            setSelectedTerritory(territoryId);
            setAttackTarget(null);
          }
          return;
        }

        // Clicking enemy territory
        if (!selectedTerritory) {
          addToast('First select one of YOUR territories to attack from', 'warning');
          return;
        }

        const sourceTerr = gameState.territories[selectedTerritory];
        if (!sourceTerr?.neighbors.includes(territoryId)) {
          addToast('Target must be adjacent to your selected territory', 'warning');
          return;
        }

        if (attackTarget === territoryId) {
          // Clicking same target again = attack with max dice (quick repeat attack)
          doAttack(selectedTerritory, territoryId);
        } else {
          // First click on this target = set it and auto-attack with max dice
          setAttackTarget(territoryId);
          doAttack(selectedTerritory, territoryId);
        }
        return;
      }

      // === FORTIFY PHASE ===
      if (gameState.phase === 'fortify') {
        if (territory.owner !== humanId) {
          addToast('You can only fortify between your own territories', 'warning');
          return;
        }

        if (!fortifySource) {
          if (territory.troops < 2) {
            addToast('Need at least 2 troops to move from here', 'warning');
            return;
          }
          setFortifySource(territoryId);
          setFortifyTarget(null);
          return;
        }

        if (fortifySource === territoryId) {
          clearSelection();
          return;
        }

        setFortifyTarget(territoryId);
      }
    },
    [gameState, aiThinking, attackPending, selectedTerritory, attackTarget, fortifySource, api, addToast, clearSelection, runAiTurns, doAttack],
  );

  const handleAttackWithDice = useCallback(
    async (dice: number) => {
      if (!selectedTerritory || !attackTarget) return;
      doAttack(selectedTerritory, attackTarget, dice);
    },
    [selectedTerritory, attackTarget, doAttack],
  );

  const handleEndPhase = useCallback(async () => {
    if (!gameState) return;
    try {
      let updated = await api.endPhase(gameState.id);
      setGameState(updated);
      clearSelection();

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

        if (gameState.freeFortify) {
          // Free fortify: let the player make additional moves; they click "End Turn" when done
          return;
        }

        updated = await api.endPhase(updated.id);
        setGameState(updated);

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

  const [showSurrenderConfirm, setShowSurrenderConfirm] = useState(false);

  const handleSurrender = useCallback(async () => {
    if (!gameState) return;
    try {
      const updated = await api.surrender(gameState.id);
      setGameState(updated);
      setShowSurrenderConfirm(false);
    } catch (e) {
      addToast(e instanceof Error ? e.message : 'Surrender failed', 'warning');
    }
  }, [gameState, api, addToast]);

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
      for (const [tid, t] of Object.entries(gameState.territories)) {
        if (t.owner === humanId && tid !== fortifySource) {
          validTargets.push(tid);
        }
      }
    }
  }

  // === RENDER ===

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center" style={{ background: '#1a1a2e' }}>
        <div className="w-6 h-6 border-2 border-red-500 border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  if (!gameState) {
    return <SetupScreen onGameStart={handleGameStart} />;
  }

  const currentPlayer = gameState.players.find((p) => p.id === gameState.currentPlayer);
  const isHumanTurn = currentPlayer && !currentPlayer.isAI;

  return (
    <div className="h-screen flex flex-col overflow-hidden" style={{ background: '#1a1a2e' }}>
      <ToastContainer toasts={toasts} onRemove={removeToast} />

      {/* Instruction bar */}
      {instructionText && isHumanTurn && !aiThinking && (
        <div
          className="px-4 py-2 text-center text-sm font-medium border-b flex-shrink-0"
          style={{
            background: 'linear-gradient(90deg, rgba(15,52,96,0.9), rgba(22,33,62,0.9))',
            borderColor: 'rgba(233,69,96,0.3)',
            color: '#e0e0e0',
          }}
        >
          {instructionText}
        </div>
      )}

      {/* AI thinking overlay */}
      {aiThinking && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
          <div
            className="px-8 py-4 rounded-xl shadow-2xl border"
            style={{ background: '#16213e', borderColor: 'rgba(233, 69, 96, 0.3)' }}
          >
            <div className="flex items-center gap-3">
              <div className="w-5 h-5 border-2 border-red-500 border-t-transparent rounded-full animate-spin" />
              <span className="text-white font-semibold">
                {currentPlayer?.name || 'AI'} is thinking...
              </span>
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
                localStorage.removeItem(GAME_ID_KEY);
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
            attackTarget={attackTarget}
            validTargets={validTargets}
            onTerritoryClick={handleTerritoryClick}
          />
        </div>

        {/* Side panel - 30% */}
        <div className="flex-[3] min-w-[280px] max-w-[380px] flex flex-col border-l border-gray-700/30">
          <div className="flex-[6] min-h-0 overflow-y-auto">
            <GameControls
              gameState={gameState}
              onEndPhase={handleEndPhase}
              onTradeCards={handleTradeCards}
              selectedTerritory={selectedTerritory}
              attackTarget={attackTarget}
              onAttack={handleAttackWithDice}
              onBlitz={handleBlitz}
              onDeployAll={handleDeployAll}
              lastPlacedTerritory={lastPlacedTerritory}
              onConquestMove={handleConquestMove}
              onFortify={handleFortify}
              fortifySource={fortifySource}
              fortifyTarget={fortifyTarget}
              onSurrender={handleSurrender}
              showSurrenderConfirm={showSurrenderConfirm}
              setShowSurrenderConfirm={setShowSurrenderConfirm}
            />
          </div>
          <div className="flex-[4] min-h-0 border-t border-gray-700/30">
            <HistoryPanel log={gameState.log} />
          </div>
        </div>
      </div>
    </div>
  );
}
