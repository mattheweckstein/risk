export type Phase = 'setup' | 'place' | 'attack' | 'fortify' | 'ended';

export interface Territory {
  id: string;
  name: string;
  continent: string;
  neighbors: string[];
  owner: string;
  troops: number;
}

export interface Player {
  id: string;
  name: string;
  isAI: boolean;
  color: string;
  cards: Card[];
  isAlive: boolean;
}

export interface Card {
  territory: string;
  type: 'infantry' | 'cavalry' | 'artillery' | 'wild';
}

export interface GameState {
  id: string;
  phase: Phase;
  turn: number;
  currentPlayer: string;
  players: Player[];
  territories: Record<string, Territory>;
  deck: Card[];
  log: LogEntry[];
  winner?: string;
  troopsToDeploy: number;
  cardTradeCount: number;
  conqueredThisTurn: boolean;
  lastAttackResult?: AttackResult;
}

export interface AttackResult {
  attackerRolls: number[];
  defenderRolls: number[];
  attackerLosses: number;
  defenderLosses: number;
  conquered: boolean;
  attackingTerritory: string;
  defendingTerritory: string;
}

export interface LogEntry {
  turn: number;
  player: string;
  message: string;
}
