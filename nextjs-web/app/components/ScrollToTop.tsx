'use client';

interface ScrollToTopProps {
  show?: boolean;
  onClick?: () => void;
}

const ScrollToTop = ({ show, onClick }: ScrollToTopProps) => {
  // 如果没有传入 show 属性，则使用组件内部状态
  const isVisible = show !== undefined ? show : false;
  
  // 处理点击事件
  const handleScrollToTop = () => {
    if (onClick) {
      onClick();
    } else {
      // 默认行为
      window.scrollTo({
        top: 0,
        behavior: 'smooth',
      });
    }
  };

  return (
    <button
      onClick={handleScrollToTop}
      className="fixed bottom-[100px] right-[20px] cursor-pointer text-[#1e80ff] dark:text-[#1e80ff] w-10 h-10 flex items-center justify-center bg-white dark:bg-gray-800 rounded-full shadow-lg transition-opacity duration-300"
      style={{ display: isVisible ? 'flex' : 'none' }}
      aria-label="Scroll to top"
    >
      <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 10l7-7m0 0l7 7m-7-7v18" />
      </svg>
    </button>
  );
};

export default ScrollToTop;