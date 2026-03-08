import { describe, it, expect } from 'vitest';
import type { GameState, Phase, Player, Territory, Card, LogEntry, AttackResult, PendingConquest } from './game';

describe('Game types', () => {
  describe('GameState interface', () => {
    it('accepts a valid GameState object', () => {
      const player: Player = {
        id: 'p1',
        name: 'Alice',
        isAI: false,
        color: 'red',
        cards: [],
        isAlive: true,
      };

      const territory: Territory = {
        id: 'alaska',
        name: 'Alaska',
        continent: 'north_america',
        neighbors: ['northwest_territory', 'kamchatka', 'alberta'],
        owner: 'p1',
        troops: 5,
      };

      const logEntry: LogEntry = {
        turn: 1,
        player: 'p1',
        message: 'Placed troops on Alaska',
      };

      const gameState: GameState = {
        id: 'game-123',
        phase: 'place',
        turn: 1,
        currentPlayer: 'p1',
        players: [player],
        territories: { alaska: territory },
        deck: [],
        log: [logEntry],
        troopsToDeploy: 3,
        cardTradeCount: 0,
        conqueredThisTurn: false,
        freeFortify: false,
      };

      expect(gameState.id).toBe('game-123');
      expect(gameState.phase).toBe('place');
      expect(gameState.players).toHaveLength(1);
      expect(gameState.territories.alaska.troops).toBe(5);
      expect(gameState.troopsToDeploy).toBe(3);
      expect(gameState.freeFortify).toBe(false);
    });

    it('accepts optional fields winner, lastAttackResult, and pendingConquest', () => {
      const attackResult: AttackResult = {
        attackerRolls: [6, 5, 3],
        defenderRolls: [4, 2],
        attackerLosses: 0,
        defenderLosses: 2,
        conquered: true,
        attackingTerritory: 'alaska',
        defendingTerritory: 'kamchatka',
      };

      const pendingConquest: PendingConquest = {
        from: 'alaska',
        to: 'kamchatka',
        minTroops: 3,
        maxTroops: 5,
      };

      const gameState: GameState = {
        id: 'game-456',
        phase: 'ended',
        turn: 10,
        currentPlayer: 'p1',
        players: [],
        territories: {},
        deck: [],
        log: [],
        winner: 'p1',
        troopsToDeploy: 0,
        cardTradeCount: 3,
        conqueredThisTurn: true,
        freeFortify: true,
        lastAttackResult: attackResult,
        pendingConquest: pendingConquest,
      };

      expect(gameState.winner).toBe('p1');
      expect(gameState.lastAttackResult?.conquered).toBe(true);
      expect(gameState.pendingConquest?.maxTroops).toBe(5);
    });
  });

  describe('Phase values', () => {
    it('all Phase values are valid strings', () => {
      const validPhases: Phase[] = ['setup', 'place', 'attack', 'fortify', 'ended'];

      expect(validPhases).toHaveLength(5);
      validPhases.forEach((phase) => {
        expect(typeof phase).toBe('string');
      });
    });

    it('covers all expected game phases', () => {
      const phases: Phase[] = ['setup', 'place', 'attack', 'fortify', 'ended'];
      expect(phases).toContain('setup');
      expect(phases).toContain('place');
      expect(phases).toContain('attack');
      expect(phases).toContain('fortify');
      expect(phases).toContain('ended');
    });
  });

  describe('Card types', () => {
    it('all card types are valid', () => {
      const cards: Card[] = [
        { territory: 'alaska', type: 'infantry' },
        { territory: 'brazil', type: 'cavalry' },
        { territory: 'china', type: 'artillery' },
        { territory: '', type: 'wild' },
      ];

      expect(cards).toHaveLength(4);
      expect(cards[0].type).toBe('infantry');
      expect(cards[1].type).toBe('cavalry');
      expect(cards[2].type).toBe('artillery');
      expect(cards[3].type).toBe('wild');
    });
  });
});
