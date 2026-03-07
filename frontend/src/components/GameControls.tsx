import { useState } from 'react';
import { GameState, Phase, AttackResult } from '../types/game';
import { continentTerritories, continentBonuses, continentDisplayNames } from '../data/continents';

interface GameControlsProps {
  gameState: GameState;
  onEndPhase: () => void;
  onTradeCards: (indices: number[]) => void;
  selectedTerritory: string | null;
  attackTarget: string | null;
  onAttack: (dice: number) => void;
  onFortify: (troops: number) => void;
  fortifySource: string | null;
  fortifyTarget: string | null;
}

const phaseLabels: Record<Phase, string> = {
  setup: 'SETUP',
  place: 'DEPLOY',
  attack: 'ATTACK',
  fortify: 'FORTIFY',
  ended: 'GAME OVER',
};

const phaseOrder: Phase[] = ['setup', 'place', 'attack', 'fortify'];

const colorMap: Record<string, string> = {
  red: '#e94560',
  blue: '#4a9eff',
  green: '#50c878',
  yellow: '#ffd700',
};

function DiceDisplay({ result }: { result: AttackResult }) {
  return (
    <div className="p-3 rounded-lg mb-3" style={{ background: 'rgba(0,0,0,0.3)' }}>
      <div className="text-xs font-semibold text-gray-400 mb-2 uppercase">Last Attack Result</div>
      <div className="flex justify-between mb-2">
        <div>
          <span className="text-xs text-red-400 mr-1">ATK:</span>
          {result.attackerRolls.map((d, i) => (
            <span
              key={i}
              className="inline-flex items-center justify-center w-8 h-8 rounded text-sm font-bold mr-1"
              style={{ background: '#e94560', color: 'white' }}
            >
              {d}
            </span>
          ))}
        </div>
        <div>
          <span className="text-xs text-blue-400 mr-1">DEF:</span>
          {result.defenderRolls.map((d, i) => (
            <span
              key={i}
              className="inline-flex items-center justify-center w-8 h-8 rounded text-sm font-bold mr-1"
              style={{ background: '#4a9eff', color: 'white' }}
            >
              {d}
            </span>
          ))}
        </div>
      </div>
      <div className="text-xs text-gray-300 font-medium">
        {result.conquered ? (
          <span className="text-yellow-400 font-bold text-sm">TERRITORY CONQUERED!</span>
        ) : (
          <>
            You lost {result.attackerLosses} | Enemy lost {result.defenderLosses}
          </>
        )}
      </div>
    </div>
  );
}

export default function GameControls({
  gameState,
  onEndPhase,
  onTradeCards,
  selectedTerritory,
  attackTarget,
  onAttack,
  onFortify,
  fortifySource,
  fortifyTarget,
}: GameControlsProps) {
  const [selectedCards, setSelectedCards] = useState<number[]>([]);
  const [fortifyTroops, setFortifyTroops] = useState(1);

  const currentPlayer = gameState.players.find((p) => p.id === gameState.currentPlayer);
  const humanPlayer = gameState.players.find((p) => !p.isAI);
  const isHumanTurn = currentPlayer && !currentPlayer.isAI;

  const continentOwners = Object.entries(continentTerritories).map(([continent, territories]) => {
    const owners = new Set(territories.map((tid) => gameState.territories[tid]?.owner).filter(Boolean));
    const owner = owners.size === 1 ? Array.from(owners)[0] : null;
    const ownerPlayer = owner ? gameState.players.find((p) => p.id === owner) : null;
    return { continent, owner: ownerPlayer, bonus: continentBonuses[continent] };
  });

  const maxAttackDice = selectedTerritory
    ? Math.min(3, (gameState.territories[selectedTerritory]?.troops || 1) - 1)
    : 0;

  const maxFortifyTroops = fortifySource
    ? (gameState.territories[fortifySource]?.troops || 1) - 1
    : 0;

  const handleCardToggle = (idx: number) => {
    setSelectedCards((prev) =>
      prev.includes(idx) ? prev.filter((i) => i !== idx) : prev.length < 3 ? [...prev, idx] : prev
    );
  };

  return (
    <div className="h-full flex flex-col overflow-y-auto p-4 gap-3" style={{ background: '#16213e' }}>
      {/* Phase indicator */}
      <div className="flex items-center gap-1">
        {phaseOrder.map((phase, i) => (
          <div key={phase} className="flex items-center">
            <div
              className={`px-2 py-1 rounded text-xs font-bold uppercase ${
                gameState.phase === phase ? 'text-white' : 'text-gray-500'
              }`}
              style={{
                background: gameState.phase === phase ? '#e94560' : 'rgba(255,255,255,0.05)',
              }}
            >
              {phaseLabels[phase]}
            </div>
            {i < phaseOrder.length - 1 && (
              <span className="text-gray-600 mx-0.5 text-xs">&#8250;</span>
            )}
          </div>
        ))}
      </div>

      {/* Current player */}
      <div
        className="p-3 rounded-lg border"
        style={{
          background: 'rgba(0,0,0,0.2)',
          borderColor: colorMap[currentPlayer?.color || ''] || '#555',
        }}
      >
        <div className="flex items-center gap-2 mb-1">
          <div
            className="w-3 h-3 rounded-full"
            style={{ background: colorMap[currentPlayer?.color || ''] || '#555' }}
          />
          <span className="font-bold text-sm text-white">
            {currentPlayer?.name || 'Unknown'}
          </span>
          {currentPlayer?.isAI && <span className="text-xs text-gray-400 ml-1">(AI)</span>}
        </div>
        <div className="text-xs text-gray-400">
          Turn {gameState.turn} | Phase: {phaseLabels[gameState.phase]}
        </div>
      </div>

      {/* Troops to deploy */}
      {(gameState.phase === 'place' || gameState.phase === 'setup') && isHumanTurn && (
        <div className="p-3 rounded-lg" style={{ background: 'rgba(233, 69, 96, 0.15)' }}>
          <div className="text-xs text-gray-400 uppercase mb-1">Troops to Deploy</div>
          <div className="text-3xl font-bold text-white">{gameState.troopsToDeploy}</div>
        </div>
      )}

      {/* Attack info panel */}
      {gameState.phase === 'attack' && isHumanTurn && (
        <div className="p-3 rounded-lg border" style={{ background: 'rgba(233,69,96,0.08)', borderColor: 'rgba(233,69,96,0.25)' }}>
          {!selectedTerritory && !attackTarget && (
            <div className="text-sm text-gray-300">
              <span className="text-red-400 font-semibold">Attack Phase:</span> Click a territory you own with 2+ troops to begin attacking.
            </div>
          )}
          {selectedTerritory && !attackTarget && (
            <div className="text-sm text-gray-300">
              <span className="text-white font-semibold">
                {gameState.territories[selectedTerritory]?.name}
              </span>
              <span className="text-gray-400"> ({gameState.territories[selectedTerritory]?.troops} troops)</span>
              <div className="text-xs text-gray-400 mt-1">Click an adjacent enemy territory to attack (yellow highlights)</div>
            </div>
          )}
          {selectedTerritory && attackTarget && (
            <>
              <div className="text-sm mb-2">
                <span className="text-white font-semibold">
                  {gameState.territories[selectedTerritory]?.name}
                </span>
                <span className="text-gray-400"> ({gameState.territories[selectedTerritory]?.troops})</span>
                <span className="text-red-400 mx-2">&#10132;</span>
                <span className="text-white font-semibold">
                  {gameState.territories[attackTarget]?.name}
                </span>
                <span className="text-gray-400"> ({gameState.territories[attackTarget]?.troops})</span>
              </div>
              <div className="text-xs text-gray-400 mb-2">Click target again to re-attack, or pick dice:</div>
              <div className="flex gap-2">
                {[1, 2, 3].map((dice) => (
                  <button
                    key={dice}
                    onClick={() => onAttack(dice)}
                    disabled={dice > maxAttackDice}
                    className="flex-1 py-2 rounded font-bold text-sm disabled:opacity-30 disabled:cursor-not-allowed hover:brightness-125 transition-all"
                    style={{
                      background: dice <= maxAttackDice ? '#e94560' : '#333',
                      color: 'white',
                    }}
                  >
                    {dice} {dice === 1 ? 'Die' : 'Dice'}
                  </button>
                ))}
              </div>
            </>
          )}
        </div>
      )}

      {/* Last attack results */}
      {gameState.lastAttackResult && <DiceDisplay result={gameState.lastAttackResult} />}

      {/* Fortify dialog */}
      {gameState.phase === 'fortify' && isHumanTurn && (
        <div className="p-3 rounded-lg border" style={{ background: 'rgba(74,158,255,0.08)', borderColor: 'rgba(74,158,255,0.25)' }}>
          {!fortifySource && (
            <div className="text-sm text-gray-300">
              <span className="text-blue-400 font-semibold">Fortify Phase:</span> Select a territory to move troops from, or skip.
            </div>
          )}
          {fortifySource && !fortifyTarget && (
            <div className="text-sm text-gray-300">
              <span className="text-white font-semibold">
                {gameState.territories[fortifySource]?.name}
              </span>
              <span className="text-gray-400"> ({gameState.territories[fortifySource]?.troops} troops)</span>
              <div className="text-xs text-gray-400 mt-1">Select a connected friendly territory to move to</div>
            </div>
          )}
          {fortifySource && fortifyTarget && (
            <>
              <div className="text-sm mb-2">
                <span className="text-white font-semibold">
                  {gameState.territories[fortifySource]?.name}
                </span>
                <span className="text-blue-400 mx-2">&#10132;</span>
                <span className="text-white font-semibold">
                  {gameState.territories[fortifyTarget]?.name}
                </span>
              </div>
              <div className="flex items-center gap-3 mb-3">
                <button
                  onClick={() => setFortifyTroops(Math.max(1, fortifyTroops - 1))}
                  className="w-8 h-8 rounded bg-blue-500/30 text-white font-bold hover:bg-blue-500/50"
                >
                  -
                </button>
                <span className="text-xl font-bold text-white min-w-[3ch] text-center">{fortifyTroops}</span>
                <button
                  onClick={() => setFortifyTroops(Math.min(maxFortifyTroops, fortifyTroops + 1))}
                  className="w-8 h-8 rounded bg-blue-500/30 text-white font-bold hover:bg-blue-500/50"
                >
                  +
                </button>
                <span className="text-xs text-gray-400">/ {maxFortifyTroops}</span>
              </div>
              <button
                onClick={() => {
                  onFortify(fortifyTroops);
                  setFortifyTroops(1);
                }}
                className="w-full py-2 rounded font-bold text-sm transition-all hover:brightness-125"
                style={{ background: '#4a9eff', color: 'white' }}
              >
                Move Troops
              </button>
            </>
          )}
        </div>
      )}

      {/* Action buttons */}
      {isHumanTurn && (
        <div className="flex flex-col gap-2">
          {gameState.phase === 'attack' && (
            <button
              onClick={onEndPhase}
              className="w-full py-2.5 rounded-lg font-bold text-sm uppercase tracking-wide transition-all hover:brightness-125"
              style={{ background: '#e94560', color: 'white' }}
            >
              End Attack Phase
            </button>
          )}
          {gameState.phase === 'fortify' && (
            <button
              onClick={onEndPhase}
              className="w-full py-2.5 rounded-lg font-bold text-sm uppercase tracking-wide transition-all hover:brightness-125"
              style={{ background: '#0f3460', color: 'white', border: '1px solid rgba(74,158,255,0.4)' }}
            >
              Skip Fortify / End Turn
            </button>
          )}
        </div>
      )}

      {/* Cards */}
      {humanPlayer && humanPlayer.cards.length > 0 && isHumanTurn && gameState.phase === 'place' && (
        <div className="p-3 rounded-lg" style={{ background: 'rgba(0,0,0,0.2)' }}>
          <div className="text-xs text-gray-400 uppercase mb-2 font-bold">Your Cards</div>
          <div className="flex flex-wrap gap-1 mb-2">
            {humanPlayer.cards.map((card, idx) => {
              const icons: Record<string, string> = {
                infantry: '\u{1F6E1}',
                cavalry: '\u{1F40E}',
                artillery: '\u{1F4A3}',
                wild: '\u{2B50}',
              };
              return (
                <button
                  key={idx}
                  onClick={() => handleCardToggle(idx)}
                  className={`px-2 py-1.5 rounded text-xs font-semibold border transition-all ${
                    selectedCards.includes(idx)
                      ? 'border-yellow-400 bg-yellow-400/20 text-yellow-300'
                      : 'border-gray-600 bg-gray-800 text-gray-300 hover:border-gray-400'
                  }`}
                >
                  {icons[card.type] || '?'} {card.type}
                </button>
              );
            })}
          </div>
          {selectedCards.length === 3 && (
            <button
              onClick={() => {
                onTradeCards(selectedCards);
                setSelectedCards([]);
              }}
              className="w-full py-2 rounded font-bold text-sm transition-all hover:brightness-125"
              style={{ background: '#ffd700', color: '#1a1a2e' }}
            >
              Trade 3 Cards
            </button>
          )}
        </div>
      )}

      {/* Players overview */}
      <div className="p-3 rounded-lg" style={{ background: 'rgba(0,0,0,0.2)' }}>
        <div className="text-xs text-gray-400 uppercase mb-2 font-bold">Players</div>
        <div className="space-y-1.5">
          {gameState.players.map((player) => {
            const terrCount = Object.values(gameState.territories).filter(
              (t) => t.owner === player.id
            ).length;
            const totalTroops = Object.values(gameState.territories)
              .filter((t) => t.owner === player.id)
              .reduce((sum, t) => sum + t.troops, 0);
            return (
              <div
                key={player.id}
                className={`flex items-center gap-2 text-xs ${
                  !player.isAlive ? 'opacity-40 line-through' : ''
                } ${gameState.currentPlayer === player.id ? 'font-bold' : ''}`}
              >
                <div
                  className="w-2.5 h-2.5 rounded-full flex-shrink-0"
                  style={{ background: colorMap[player.color] || '#555' }}
                />
                <span className="flex-1 truncate text-gray-200">{player.name}</span>
                <span className="text-gray-500">{terrCount}T</span>
                <span className="text-gray-500">{totalTroops}A</span>
              </div>
            );
          })}
        </div>
      </div>

      {/* Continent ownership */}
      <div className="p-3 rounded-lg" style={{ background: 'rgba(0,0,0,0.2)' }}>
        <div className="text-xs text-gray-400 uppercase mb-2 font-bold">Continents</div>
        <div className="space-y-1">
          {continentOwners.map(({ continent, owner, bonus }) => (
            <div key={continent} className="flex items-center gap-2 text-xs">
              {owner ? (
                <div
                  className="w-2.5 h-2.5 rounded-full flex-shrink-0"
                  style={{ background: colorMap[owner.color] || '#555' }}
                />
              ) : (
                <div className="w-2.5 h-2.5 rounded-full flex-shrink-0 bg-gray-600" />
              )}
              <span className="flex-1 text-gray-300">{continentDisplayNames[continent]}</span>
              <span className="text-gray-500">+{bonus}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
