'use client';

import React from 'react';
import { Button, Space, DatePicker, ConfigProvider, message } from 'antd';
import { SearchOutlined, DownloadOutlined } from '@ant-design/icons';
import zhCN from 'antd/locale/zh_CN';

const AntdTest: React.FC = () => {
  const [messageApi, contextHolder] = message.useMessage();

  const showMessage = () => {
    messageApi.success('成功集成了 Ant Design 组件！');
  };

  return (
    <ConfigProvider locale={zhCN} theme={{
      token: {
        colorPrimary: '#1e80ff',
      },
    }}>
      {contextHolder}
      <div className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-bold mb-4 dark:text-white">Ant Design 组件测试</h2>
        <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
          <Space wrap>
            <Button type="primary" icon={<SearchOutlined />} onClick={showMessage}>
              测试组件
            </Button>
            <Button type="primary" icon={<DownloadOutlined />}>
              下载
            </Button>
            <Button>默认按钮</Button>
            <Button type="dashed">虚线按钮</Button>
            <Button type="text">文本按钮</Button>
            <Button type="link">链接按钮</Button>
          </Space>
          <DatePicker placeholder="选择日期" />
        </Space>
      </div>
    </ConfigProvider>
  );
};

export default AntdTest;