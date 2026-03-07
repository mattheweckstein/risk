import { useState } from 'react';
import { Territory } from '../types/game';

interface TerritoryNodeProps {
  territory: Territory;
  pathData: string;
  labelX: number;
  labelY: number;
  isSelected: boolean;
  isValidTarget: boolean;
  onClick: () => void;
  playerColor: string;
}

// Map player color names to hex values
const colorMap: Record<string, string> = {
  red: '#e94560',
  blue: '#4a9eff',
  green: '#50c878',
  yellow: '#ffd700',
};

const darkerColorMap: Record<string, string> = {
  red: '#b8354d',
  blue: '#3a7ecc',
  green: '#40a060',
  yellow: '#ccab00',
};

export default function TerritoryNode({
  territory,
  pathData,
  labelX,
  labelY,
  isSelected,
  isValidTarget,
  onClick,
  playerColor,
}: TerritoryNodeProps) {
  const [hovered, setHovered] = useState(false);

  const fillColor = colorMap[playerColor] || playerColor || '#555';
  const darkerFill = darkerColorMap[playerColor] || fillColor;

  // Selected territory gets bright highlight
  const strokeColor = isSelected
    ? '#ffffff'
    : isValidTarget
    ? '#ffdd57'
    : hovered
    ? '#ffffff88'
    : '#1a1a2e';

  const strokeWidth = isSelected ? 3 : isValidTarget ? 2.5 : hovered ? 1.5 : 0.8;

  const fillOpacity = territory.owner ? 0.85 : 0.4;

  return (
    <g
      onClick={onClick}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
      className="territory-path"
      style={{ cursor: 'pointer' }}
    >
      {/* Territory shape */}
      <path
        d={pathData}
        fill={fillColor}
        fillOpacity={fillOpacity}
        stroke={strokeColor}
        strokeWidth={strokeWidth}
        className={isValidTarget ? 'animate-pulse-target' : ''}
        style={{
          transition: 'fill-opacity 0.2s, stroke 0.2s',
          filter: isSelected ? 'brightness(1.3) drop-shadow(0 0 8px rgba(255,255,255,0.4))' : undefined,
        }}
      />

      {/* Troop count badge */}
      <circle
        cx={labelX}
        cy={labelY}
        r={territory.troops >= 10 ? 13 : 11}
        fill={darkerFill}
        stroke="#0d0d1a"
        strokeWidth={1.5}
        style={{ filter: 'drop-shadow(0 1px 3px rgba(0,0,0,0.5))' }}
      />
      <text
        x={labelX}
        y={labelY}
        textAnchor="middle"
        dominantBaseline="central"
        fill="white"
        fontSize={territory.troops >= 10 ? 10 : 11}
        fontWeight="bold"
        style={{ pointerEvents: 'none', userSelect: 'none' }}
      >
        {territory.troops}
      </text>

      {/* Tooltip on hover */}
      {hovered && (
        <g style={{ pointerEvents: 'none' }}>
          <rect
            x={labelX - 65}
            y={labelY - 52}
            width={130}
            height={38}
            rx={6}
            fill="rgba(15, 52, 96, 0.95)"
            stroke="rgba(233, 69, 96, 0.5)"
            strokeWidth={1}
          />
          <text
            x={labelX}
            y={labelY - 38}
            textAnchor="middle"
            fill="white"
            fontSize={9}
            fontWeight="bold"
          >
            {territory.name}
          </text>
          <text
            x={labelX}
            y={labelY - 25}
            textAnchor="middle"
            fill="#aaa"
            fontSize={8}
          >
            Troops: {territory.troops} | {territory.continent.replace('_', ' ')}
          </text>
        </g>
      )}
    </g>
  );
}
