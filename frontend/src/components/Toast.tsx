import { useEffect, useState } from 'react';

export interface ToastMessage {
  id: number;
  text: string;
  type: 'info' | 'success' | 'warning';
}

interface ToastContainerProps {
  toasts: ToastMessage[];
  onRemove: (id: number) => void;
}

function ToastItem({ toast, onRemove }: { toast: ToastMessage; onRemove: () => void }) {
  const [fading, setFading] = useState(false);

  useEffect(() => {
    const fadeTimer = setTimeout(() => setFading(true), 3000);
    const removeTimer = setTimeout(onRemove, 3300);
    return () => {
      clearTimeout(fadeTimer);
      clearTimeout(removeTimer);
    };
  }, [onRemove]);

  const borderColors = {
    info: 'rgba(74, 158, 255, 0.6)',
    success: 'rgba(80, 200, 120, 0.6)',
    warning: 'rgba(255, 215, 0, 0.6)',
  };

  return (
    <div
      className={`toast ${fading ? 'fade-out' : ''}`}
      style={{ borderColor: borderColors[toast.type] }}
    >
      {toast.text}
    </div>
  );
}

export default function ToastContainer({ toasts, onRemove }: ToastContainerProps) {
  if (toasts.length === 0) return null;

  return (
    <div className="toast-container">
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onRemove={() => onRemove(toast.id)} />
      ))}
    </div>
  );
}
