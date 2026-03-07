import { GameState } from '../types/game';

async function apiFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `API error: ${res.status}`);
  }
  return res.json();
}

export function useGameApi() {
  const createGame = (playerName: string, aiCount: number, aiNames: string[], freeFortify: boolean): Promise<GameState> =>
    apiFetch<GameState>('/api/game/new', {
      method: 'POST',
      body: JSON.stringify({ playerName, aiCount, aiNames, freeFortify }),
    });

  const getGame = (id: string): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}`);

  const placeTroops = (id: string, territory: string, troops: number): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/place`, {
      method: 'POST',
      body: JSON.stringify({ territory, troops }),
    });

  const attack = (id: string, from: string, to: string, attackerDice: number): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/attack`, {
      method: 'POST',
      body: JSON.stringify({ from, to, attackerDice }),
    });

  const fortify = (id: string, from: string, to: string, troops: number): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/fortify`, {
      method: 'POST',
      body: JSON.stringify({ from, to, troops }),
    });

  const endPhase = (id: string): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/end-phase`, {
      method: 'POST',
    });

  const tradeCards = (id: string, cardIndices: number[]): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/cards/trade`, {
      method: 'POST',
      body: JSON.stringify({ cardIndices }),
    });

  const moveAfterConquest = (id: string, troops: number): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/attack/move`, {
      method: 'POST',
      body: JSON.stringify({ troops }),
    });

  const aiTurn = (id: string): Promise<GameState> =>
    apiFetch<GameState>(`/api/game/${id}/ai-turn`);

  return { createGame, getGame, placeTroops, attack, moveAfterConquest, fortify, endPhase, tradeCards, aiTurn };
}
