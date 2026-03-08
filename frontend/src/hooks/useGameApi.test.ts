import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useGameApi } from './useGameApi';

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

function mockJsonResponse(data: unknown, ok = true, status = 200) {
  return Promise.resolve({
    ok,
    status,
    json: () => Promise.resolve(data),
    text: () => Promise.resolve(JSON.stringify(data)),
  });
}

function mockErrorResponse(message: string, status = 400) {
  return Promise.resolve({
    ok: false,
    status,
    text: () => Promise.resolve(message),
  });
}

describe('useGameApi', () => {
  let api: ReturnType<typeof useGameApi>;

  beforeEach(() => {
    mockFetch.mockReset();
    api = useGameApi();
  });

  describe('createGame', () => {
    it('sends correct request body and returns game state', async () => {
      const mockGame = { id: 'game-1', phase: 'setup', players: [] };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      const result = await api.createGame('Alice', 2, ['Bot1', 'Bot2'], true);

      expect(mockFetch).toHaveBeenCalledWith('/api/game/new', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ playerName: 'Alice', aiCount: 2, aiNames: ['Bot1', 'Bot2'], freeFortify: true }),
      });
      expect(result).toEqual(mockGame);
    });
  });

  describe('getGame', () => {
    it('sends GET request with correct game ID', async () => {
      const mockGame = { id: 'game-1', phase: 'place' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      const result = await api.getGame('game-1');

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1', {
        headers: { 'Content-Type': 'application/json' },
      });
      expect(result).toEqual(mockGame);
    });
  });

  describe('placeTroops', () => {
    it('sends correct params', async () => {
      const mockGame = { id: 'game-1', phase: 'place' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.placeTroops('game-1', 'alaska', 3);

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/place', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ territory: 'alaska', troops: 3 }),
      });
    });
  });

  describe('attack', () => {
    it('sends correct params', async () => {
      const mockGame = { id: 'game-1', phase: 'attack' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.attack('game-1', 'alaska', 'kamchatka', 3);

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/attack', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ from: 'alaska', to: 'kamchatka', attackerDice: 3 }),
      });
    });
  });

  describe('fortify', () => {
    it('sends correct params', async () => {
      const mockGame = { id: 'game-1', phase: 'fortify' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.fortify('game-1', 'alaska', 'northwest_territory', 2);

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/fortify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ from: 'alaska', to: 'northwest_territory', troops: 2 }),
      });
    });
  });

  describe('endPhase', () => {
    it('sends POST request', async () => {
      const mockGame = { id: 'game-1', phase: 'fortify' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.endPhase('game-1');

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/end-phase', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
    });
  });

  describe('tradeCards', () => {
    it('sends card indices', async () => {
      const mockGame = { id: 'game-1', phase: 'place' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.tradeCards('game-1', [0, 1, 2]);

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/cards/trade', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ cardIndices: [0, 1, 2] }),
      });
    });
  });

  describe('moveAfterConquest', () => {
    it('sends troops count', async () => {
      const mockGame = { id: 'game-1', phase: 'attack' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.moveAfterConquest('game-1', 5);

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/attack/move', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ troops: 5 }),
      });
    });
  });

  describe('aiTurn', () => {
    it('uses GET method (default, no method specified)', async () => {
      const mockGame = { id: 'game-1', phase: 'place' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.aiTurn('game-1');

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/ai-turn', {
        headers: { 'Content-Type': 'application/json' },
      });
      // Verify no method was explicitly set (defaults to GET)
      const callArgs = mockFetch.mock.calls[0][1];
      expect(callArgs.method).toBeUndefined();
    });
  });

  describe('surrender', () => {
    it('sends POST request', async () => {
      const mockGame = { id: 'game-1', phase: 'ended' };
      mockFetch.mockReturnValue(mockJsonResponse(mockGame));

      await api.surrender('game-1');

      expect(mockFetch).toHaveBeenCalledWith('/api/game/game-1/surrender', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
      });
    });
  });

  describe('error handling', () => {
    it('throws error with message from non-ok response', async () => {
      mockFetch.mockReturnValue(mockErrorResponse('Territory not owned by player'));

      await expect(api.placeTroops('game-1', 'alaska', 1)).rejects.toThrow(
        'Territory not owned by player'
      );
    });

    it('throws generic error when response body is empty', async () => {
      mockFetch.mockReturnValue(
        Promise.resolve({
          ok: false,
          status: 500,
          text: () => Promise.resolve(''),
        })
      );

      await expect(api.placeTroops('game-1', 'alaska', 1)).rejects.toThrow(
        'API error: 500'
      );
    });
  });
});
