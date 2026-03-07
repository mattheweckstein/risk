import { useState } from 'react';
import { GameState } from '../types/game';
import { useGameApi } from '../hooks/useGameApi';

interface SetupScreenProps {
  onGameStart: (game: GameState) => void;
}

const defaultAINames = ['General Rex', 'Commander Voss', 'Admiral Kane'];

export default function SetupScreen({ onGameStart }: SetupScreenProps) {
  const [playerName, setPlayerName] = useState('');
  const [aiCount, setAICount] = useState(2);
  const [aiNames, setAINames] = useState([...defaultAINames]);
  const [freeFortify, setFreeFortify] = useState(true);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const api = useGameApi();

  const handleStart = async () => {
    if (!playerName.trim()) {
      setError('Please enter your name');
      return;
    }
    setLoading(true);
    setError('');
    try {
      const game = await api.createGame(
        playerName.trim(),
        aiCount,
        aiNames.slice(0, aiCount),
        freeFortify
      );
      onGameStart(game);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to create game');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-lg">
        {/* Title */}
        <div className="text-center mb-10">
          <h1 className="text-6xl font-bold tracking-wider mb-2" style={{ color: '#e94560' }}>
            RISK
          </h1>
          <p className="text-lg text-gray-400 tracking-widest uppercase">
            World Domination
          </p>
          <div className="mt-4 h-0.5 bg-gradient-to-r from-transparent via-red-500 to-transparent"></div>
        </div>

        {/* Form Card */}
        <div
          className="rounded-xl p-8 shadow-2xl border"
          style={{
            backgroundColor: '#16213e',
            borderColor: 'rgba(233, 69, 96, 0.2)',
          }}
        >
          {/* Player Name */}
          <div className="mb-6">
            <label className="block text-sm font-semibold mb-2 text-gray-300 uppercase tracking-wide">
              Your Name
            </label>
            <input
              type="text"
              value={playerName}
              onChange={(e) => setPlayerName(e.target.value)}
              placeholder="Enter your commander name..."
              className="w-full px-4 py-3 rounded-lg border bg-[#1a1a2e] text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all"
              style={{ borderColor: 'rgba(233, 69, 96, 0.3)' }}
              onKeyDown={(e) => e.key === 'Enter' && handleStart()}
            />
          </div>

          {/* AI Count */}
          <div className="mb-6">
            <label className="block text-sm font-semibold mb-3 text-gray-300 uppercase tracking-wide">
              AI Opponents
            </label>
            <div className="flex gap-3">
              {[1, 2, 3].map((n) => (
                <label
                  key={n}
                  className={`flex-1 text-center py-3 rounded-lg cursor-pointer border-2 transition-all font-semibold ${
                    aiCount === n
                      ? 'border-red-500 bg-red-500/20 text-white'
                      : 'border-gray-600 bg-[#1a1a2e] text-gray-400 hover:border-gray-500'
                  }`}
                >
                  <input
                    type="radio"
                    name="aiCount"
                    value={n}
                    checked={aiCount === n}
                    onChange={() => setAICount(n)}
                    className="sr-only"
                  />
                  {n} {n === 1 ? 'Opponent' : 'Opponents'}
                </label>
              ))}
            </div>
          </div>

          {/* AI Names */}
          <div className="mb-8">
            <label className="block text-sm font-semibold mb-3 text-gray-300 uppercase tracking-wide">
              AI Commander Names
            </label>
            <div className="space-y-2">
              {Array.from({ length: aiCount }).map((_, i) => (
                <input
                  key={i}
                  type="text"
                  value={aiNames[i]}
                  onChange={(e) => {
                    const updated = [...aiNames];
                    updated[i] = e.target.value;
                    setAINames(updated);
                  }}
                  placeholder={`AI ${i + 1} name...`}
                  className="w-full px-4 py-2.5 rounded-lg border bg-[#1a1a2e] text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all text-sm"
                  style={{ borderColor: 'rgba(233, 69, 96, 0.2)' }}
                />
              ))}
            </div>
          </div>

          {/* Game Settings */}
          <div className="mb-8">
            <label className="block text-sm font-semibold mb-3 text-gray-300 uppercase tracking-wide">
              Game Settings
            </label>
            <label
              className={`flex items-center justify-between p-3 rounded-lg border-2 cursor-pointer transition-all ${
                freeFortify
                  ? 'border-red-500 bg-red-500/10'
                  : 'border-gray-600 bg-[#1a1a2e] hover:border-gray-500'
              }`}
              onClick={() => setFreeFortify(!freeFortify)}
            >
              <div>
                <div className="text-sm font-semibold text-white">Free Fortify</div>
                <div className="text-xs text-gray-400 mt-0.5">Allow multiple troop movements per turn</div>
              </div>
              <div
                className={`w-10 h-6 rounded-full relative transition-all flex-shrink-0 ml-3 ${
                  freeFortify ? 'bg-red-500' : 'bg-gray-600'
                }`}
              >
                <div
                  className={`w-4 h-4 rounded-full bg-white absolute top-1 transition-all ${
                    freeFortify ? 'left-5' : 'left-1'
                  }`}
                />
              </div>
            </label>
          </div>

          {/* Error */}
          {error && (
            <div className="mb-4 p-3 rounded-lg bg-red-500/20 border border-red-500/40 text-red-300 text-sm">
              {error}
            </div>
          )}

          {/* Start Button */}
          <button
            onClick={handleStart}
            disabled={loading}
            className="w-full py-4 rounded-lg font-bold text-lg uppercase tracking-wider transition-all disabled:opacity-50 disabled:cursor-not-allowed hover:scale-[1.02] active:scale-[0.98]"
            style={{
              background: 'linear-gradient(135deg, #e94560, #c23152)',
              color: 'white',
              boxShadow: '0 4px 15px rgba(233, 69, 96, 0.4)',
            }}
          >
            {loading ? 'Deploying Forces...' : 'Start Game'}
          </button>
        </div>

        {/* Footer */}
        <p className="text-center text-gray-600 text-xs mt-6">
          Conquer the world. Eliminate all opponents. Glory awaits.
        </p>
      </div>
    </div>
  );
}
