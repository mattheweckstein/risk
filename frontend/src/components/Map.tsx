import { GameState, Territory } from '../types/game';
import { territoryPaths, connectionLines, continentOutlines } from '../data/territoryPaths';
import TerritoryNode from './TerritoryNode';

interface MapProps {
  gameState: GameState;
  selectedTerritory: string | null;
  attackTarget: string | null;
  validTargets: string[];
  onTerritoryClick: (territoryId: string) => void;
}

function getPlayerColor(gameState: GameState, territory: Territory): string {
  if (!territory.owner) return '#555';
  const player = gameState.players.find((p) => p.id === territory.owner);
  return player?.color || '#555';
}

function getTerritoryCenter(id: string): { x: number; y: number } | null {
  const tp = territoryPaths.find((t) => t.id === id);
  return tp ? { x: tp.labelX, y: tp.labelY } : null;
}

// Brighter continent colors for background outlines
const continentFillColors: Record<string, string> = {
  north_america: 'rgba(255, 200, 50, 0.06)',
  south_america: 'rgba(233, 69, 96, 0.06)',
  europe: 'rgba(100, 149, 237, 0.06)',
  africa: 'rgba(255, 165, 0, 0.06)',
  asia: 'rgba(50, 205, 50, 0.06)',
  australia: 'rgba(148, 103, 189, 0.06)',
};

const continentStrokeColors: Record<string, string> = {
  north_america: 'rgba(255, 200, 50, 0.15)',
  south_america: 'rgba(233, 69, 96, 0.15)',
  europe: 'rgba(100, 149, 237, 0.15)',
  africa: 'rgba(255, 165, 0, 0.15)',
  asia: 'rgba(50, 205, 50, 0.15)',
  australia: 'rgba(148, 103, 189, 0.15)',
};

export default function Map({
  gameState,
  selectedTerritory,
  attackTarget,
  validTargets,
  onTerritoryClick,
}: MapProps) {
  // Draw attack arrow from selected territory to attack target
  const attackArrow =
    selectedTerritory && attackTarget
      ? (() => {
          const from = getTerritoryCenter(selectedTerritory);
          const to = getTerritoryCenter(attackTarget);
          if (!from || !to) return null;
          return { from, to };
        })()
      : null;

  // Compute neighbor connection lines when a territory is selected during attack phase
  const neighborLines: { from: { x: number; y: number }; to: { x: number; y: number }; isFriendly: boolean; targetId: string }[] = [];
  if (selectedTerritory && gameState.phase === 'attack' && !attackTarget) {
    const source = gameState.territories[selectedTerritory];
    if (source) {
      const fromCenter = getTerritoryCenter(selectedTerritory);
      if (fromCenter) {
        for (const nid of source.neighbors) {
          const neighbor = gameState.territories[nid];
          if (!neighbor) continue;
          const toCenter = getTerritoryCenter(nid);
          if (!toCenter) continue;
          const isFriendly = neighbor.owner === source.owner;
          neighborLines.push({ from: fromCenter, to: toCenter, isFriendly, targetId: nid });
        }
      }
    }
  }

  return (
    <div className="relative w-full h-full">
      <svg
        viewBox="0 0 1100 600"
        className="w-full h-full"
        style={{ maxHeight: 'calc(100vh - 64px)' }}
        preserveAspectRatio="xMidYMid meet"
      >
        {/* Background */}
        <defs>
          <radialGradient id="ocean-gradient" cx="50%" cy="50%" r="60%">
            <stop offset="0%" stopColor="#1a2744" />
            <stop offset="100%" stopColor="#0d1525" />
          </radialGradient>
          <marker
            id="attack-arrow"
            viewBox="0 0 10 10"
            refX="9"
            refY="5"
            markerWidth="6"
            markerHeight="6"
            orient="auto-start-reverse"
          >
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#ff4444" />
          </marker>
        </defs>
        <rect width="1100" height="600" fill="url(#ocean-gradient)" />

        {/* Grid lines for ocean texture */}
        <g opacity={0.05}>
          {Array.from({ length: 22 }).map((_, i) => (
            <line
              key={`h${i}`}
              x1={0}
              y1={i * 30}
              x2={1100}
              y2={i * 30}
              stroke="#4a9eff"
              strokeWidth={0.5}
            />
          ))}
          {Array.from({ length: 37 }).map((_, i) => (
            <line
              key={`v${i}`}
              x1={i * 30}
              y1={0}
              x2={i * 30}
              y2={600}
              stroke="#4a9eff"
              strokeWidth={0.5}
            />
          ))}
        </g>

        {/* Continent background shapes */}
        {Object.entries(continentOutlines).map(([continent, outlinePath]) => (
          <path
            key={`continent-bg-${continent}`}
            d={outlinePath}
            fill={continentFillColors[continent] || 'rgba(255,255,255,0.03)'}
            stroke={continentStrokeColors[continent] || 'rgba(255,255,255,0.08)'}
            strokeWidth={1.5}
            strokeDasharray="6 4"
            style={{ pointerEvents: 'none' }}
          />
        ))}

        {/* Continent labels */}
        <text x={170} y={16} fill="rgba(255,200,50,0.25)" fontSize={14} fontWeight="bold" textAnchor="middle">
          NORTH AMERICA
        </text>
        <text x={210} y={490} fill="rgba(233,69,96,0.25)" fontSize={12} fontWeight="bold" textAnchor="middle">
          SOUTH AMERICA
        </text>
        <text x={540} y={45} fill="rgba(100,149,237,0.25)" fontSize={12} fontWeight="bold" textAnchor="middle">
          EUROPE
        </text>
        <text x={530} y={470} fill="rgba(255,165,0,0.25)" fontSize={12} fontWeight="bold" textAnchor="middle">
          AFRICA
        </text>
        <text x={860} y={16} fill="rgba(50,205,50,0.25)" fontSize={14} fontWeight="bold" textAnchor="middle">
          ASIA
        </text>
        <text x={975} y={450} fill="rgba(148,103,189,0.25)" fontSize={12} fontWeight="bold" textAnchor="middle">
          AUSTRALIA
        </text>

        {/* Connection lines between non-adjacent territories (cross-ocean) */}
        {connectionLines.map((conn) => (
          <path
            key={`${conn.from}-${conn.to}`}
            d={conn.path}
            fill="none"
            stroke="rgba(255,255,255,0.08)"
            strokeWidth={1}
            strokeDasharray="4 4"
          />
        ))}

        {/* Neighbor connection lines when territory selected during attack */}
        {neighborLines.map((line) => (
          <line
            key={`neighbor-${line.targetId}`}
            x1={line.from.x}
            y1={line.from.y}
            x2={line.to.x}
            y2={line.to.y}
            stroke={line.isFriendly ? 'rgba(255,255,255,0.25)' : 'rgba(255,221,87,0.4)'}
            strokeWidth={line.isFriendly ? 1 : 1.5}
            strokeDasharray={line.isFriendly ? '3 3' : '5 3'}
            style={{ pointerEvents: 'none' }}
          />
        ))}

        {/* Territories */}
        {territoryPaths.map((tp) => {
          const territory = gameState.territories[tp.id];
          if (!territory) return null;
          return (
            <TerritoryNode
              key={tp.id}
              territory={territory}
              pathData={tp.path}
              labelX={tp.labelX}
              labelY={tp.labelY}
              isSelected={selectedTerritory === tp.id}
              isAttackTarget={attackTarget === tp.id}
              isValidTarget={validTargets.includes(tp.id)}
              onClick={() => onTerritoryClick(tp.id)}
              playerColor={getPlayerColor(gameState, territory)}
            />
          );
        })}

        {/* Attack arrow */}
        {attackArrow && (
          <line
            x1={attackArrow.from.x}
            y1={attackArrow.from.y}
            x2={attackArrow.to.x}
            y2={attackArrow.to.y}
            stroke="#ff4444"
            strokeWidth={3}
            strokeDasharray="8 4"
            markerEnd="url(#attack-arrow)"
            opacity={0.8}
            style={{ pointerEvents: 'none' }}
          >
            <animate
              attributeName="stroke-dashoffset"
              values="0;-24"
              dur="0.8s"
              repeatCount="indefinite"
            />
          </line>
        )}
      </svg>
    </div>
  );
}
