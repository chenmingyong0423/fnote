'use client';

import { useState, useEffect } from 'react';

interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error' | 'info' | 'warning';
}

// 创建全局toast事件系统
let toastListeners: ((toast: Toast) => void)[] = [];

export const showToast = (message: string, type: 'success' | 'error' | 'info' | 'warning' = 'info') => {
  const newToast: Toast = {
    id: Date.now(),
    message,
    type,
  };
  
  toastListeners.forEach(listener => listener(newToast));
};

const MyToast = () => {
  const [toasts, setToasts] = useState<Toast[]>([]);

  useEffect(() => {
    const addToast = (toast: Toast) => {
      setToasts(prev => [...prev, toast]);
      
      // 自动移除toast
      setTimeout(() => {
        setToasts(prev => prev.filter(t => t.id !== toast.id));
      }, 3000);
    };
    
    toastListeners.push(addToast);
    
    return () => {
      toastListeners = toastListeners.filter(listener => listener !== addToast);
    };
  }, []);

  const removeToast = (id: number) => {
    setToasts(prev => prev.filter(toast => toast.id !== id));
  };

  const getToastClasses = (type: string) => {
    const baseClasses = 'p-4 rounded shadow-lg';
    switch (type) {
      case 'success':
        return `${baseClasses} bg-green-500 text-white`;
      case 'error':
        return `${baseClasses} bg-red-500 text-white`;
      case 'warning':
        return `${baseClasses} bg-yellow-500 text-white`;
      case 'info':
      default:
        return `${baseClasses} bg-blue-500 text-white`;
    }
  };

  return (
    <div className="fixed top-5 right-5 z-[9999] flex flex-col gap-2 max-w-sm">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className={`${getToastClasses(toast.type)} flex justify-between items-center animate-fade-in-down`}
          style={{animation: 'slideInRight 0.3s ease-out forwards'}}
        >
          <span>{toast.message}</span>
          <button
            onClick={() => removeToast(toast.id)}
            className="ml-4 text-white hover:text-gray-200"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </button>
        </div>
      ))}
    </div>
  );
};

export default MyToast;