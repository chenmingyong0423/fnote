"use client";
import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import rehypeHighlight from "rehype-highlight";
import { extractHeadingIdsByLine, genHeadingId } from "./Toc";

interface MarkdownPreviewProps {
    content: string;
    className?: string;
    headingIdGenerator?: (text: string) => string;
}

export const MarkdownPreview: React.FC<MarkdownPreviewProps> = ({content, className = ''}) => {
    const headingIdsByLine = React.useMemo(
        () => extractHeadingIdsByLine(content),
        [content]
    );
    const getHeadingId = (children: React.ReactNode, node: unknown) => {
        const line = (node as { position?: { start?: { line?: number } } })?.position?.start?.line;
        return (line ? headingIdsByLine.get(line) : undefined) || genHeadingId(String(children));
    };

    return (
        <div className={`markdown-body ${className}`}>
            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeRaw, rehypeHighlight]}
                components={{
                    h1: ({children, node, ...props}) => <h1 id={getHeadingId(children, node)} {...props}>{children}</h1>,
                    h2: ({children, node, ...props}) => <h2 id={getHeadingId(children, node)} {...props}>{children}</h2>,
                    h3: ({children, node, ...props}) => <h3 id={getHeadingId(children, node)} {...props}>{children}</h3>,
                    h4: ({children, node, ...props}) => <h4 id={getHeadingId(children, node)} {...props}>{children}</h4>,
                    h5: ({children, node, ...props}) => <h5 id={getHeadingId(children, node)} {...props}>{children}</h5>,
                    h6: ({children, node, ...props}) => <h6 id={getHeadingId(children, node)} {...props}>{children}</h6>,
                }}
            >
                {content}
            </ReactMarkdown>
        </div>
    );
};
