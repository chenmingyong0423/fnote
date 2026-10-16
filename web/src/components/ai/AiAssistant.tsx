"use client";

import { useRef, useState } from "react";
import dynamic from "next/dynamic";
import Image from "next/image";
import launcherIcon from "@/public/ai-assistant.png";
import styles from "./AiAssistant.module.css";

const AiChatPanel = dynamic(() => import("./AiChatPanel"), { ssr: false });

export default function AiAssistant({ logoSrc, assistantName = "AI 助手小勇" }: { logoSrc: string; assistantName?: string }) {
  const [open, setOpen] = useState(false);
  const [mounted, setMounted] = useState(false);
  const trigger = useRef<HTMLButtonElement>(null);

  return (
    <>
      <button
        ref={trigger}
        type="button"
        className={styles.launcher}
        aria-label={`打开 ${assistantName}`}
        aria-haspopup="dialog"
        aria-expanded={open}
        title={assistantName}
        onClick={() => { setMounted(true); setOpen(true); }}
      >
        <Image
          src={launcherIcon}
          alt=""
          aria-hidden
          width={56}
          height={56}
          sizes="56px"
          className={styles.launcherIcon}
        />
      </button>
      {mounted && (
        <AiChatPanel
          logoSrc={logoSrc}
          assistantName={assistantName}
          open={open}
          onClose={() => setOpen(false)}
          afterClose={() => trigger.current?.focus({ preventScroll: true })}
        />
      )}
    </>
  );
}
