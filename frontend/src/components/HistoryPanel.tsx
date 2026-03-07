import { useEffect, useRef } from 'react';
import { LogEntry } from '../types/game';

interface HistoryPanelProps {
  log: LogEntry[];
}

export default function HistoryPanel({ log }: HistoryPanelProps) {
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [log.length]);

  const recentLog = log.slice(-30);

  return (
    <div
      className="flex flex-col h-full"
      style={{ background: '#16213e' }}
    >
      <div className="px-4 py-2 border-b border-gray-700/50">
        <h3 className="text-xs font-bold text-gray-400 uppercase tracking-wider">
          Battle Log
        </h3>
      </div>
      <div
        ref={scrollRef}
        className="flex-1 overflow-y-auto px-4 py-2 space-y-1"
      >
        {recentLog.map((entry, idx) => (
          <div
            key={idx}
            className="text-xs leading-relaxed py-1 border-b border-gray-700/20"
          >
            <span className="text-gray-500 font-mono mr-1.5">
              T{entry.turn}
            </span>
            <span className="text-gray-300">{entry.message}</span>
          </div>
        ))}
        {recentLog.length === 0 && (
          <div className="text-xs text-gray-500 italic py-4 text-center">
            No events yet...
          </div>
        )}
      </div>
    </div>
  );
}
