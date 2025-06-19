"use client";
import { useEffect } from "react";
import { useConfigStore } from "../store/config";
import type { IndexConfigVO } from "../api/config";

export default function ConfigToZustand({ config }: { config: IndexConfigVO }) {
  const setConfig = useConfigStore((s) => s.setConfig);
  useEffect(() => {
    setConfig(config);
  }, [config, setConfig]);
  return null;
}
