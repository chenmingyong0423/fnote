"use client";
import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import rehypeHighlight from "rehype-highlight";

interface MarkdownPreviewProps {
    content: string;
    className?: string;
    headingIdGenerator?: (text: string) => string;
}

export const MarkdownPreview: React.FC<MarkdownPreviewProps> = ({content, className = ''}) => {
    return (
        <div className={`markdown-body ${className}`}>
            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeRaw, rehypeHighlight]}
                components={{
                    h1: ({children, ...props}) => <h1 id={genHeadingId(String(children))} {...props}>{children}</h1>,
                    h2: ({children, ...props}) => <h2 id={genHeadingId(String(children))} {...props}>{children}</h2>,
                    h3: ({children, ...props}) => <h3 id={genHeadingId(String(children))} {...props}>{children}</h3>,
                    h4: ({children, ...props}) => <h4 id={genHeadingId(String(children))} {...props}>{children}</h4>,
                    h5: ({children, ...props}) => <h5 id={genHeadingId(String(children))} {...props}>{children}</h5>,
                    h6: ({children, ...props}) => <h6 id={genHeadingId(String(children))} {...props}>{children}</h6>,
                }}
            >
                {content}
            </ReactMarkdown>
        </div>
    );
};

export function genHeadingId(text: string) {
    return text.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, "-").replace(/^-+|-+$/g, "");
}
