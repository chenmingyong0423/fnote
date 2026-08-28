"use client";
import React from "react";
import ReactMarkdown from "react-markdown";
import type { Components } from "react-markdown";
import { CheckOutlined, CopyOutlined } from "@ant-design/icons";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import rehypeHighlight from "rehype-highlight";
import { extractHeadingIdsByLine, genHeadingId } from "./Toc";

interface MarkdownPreviewProps {
    content: string;
    className?: string;
    headingIdGenerator?: (text: string) => string;
    headingLevelOffset?: number;
    theme?: "default" | "blog";
    signatureText?: string;
}

interface BlogCodeBlockProps extends React.HTMLAttributes<HTMLPreElement> {
    children: React.ReactNode;
    signatureText: string;
}

const getTextContent = (node: React.ReactNode): string => {
    if (typeof node === "string" || typeof node === "number") {
        return String(node);
    }

    if (Array.isArray(node)) {
        return node.map(getTextContent).join("");
    }

    if (React.isValidElement<{ children?: React.ReactNode }>(node)) {
        return getTextContent(node.props.children);
    }

    return "";
};

const getCodeLanguage = (children: React.ReactNode) => {
    const child = React.Children.toArray(children).find(React.isValidElement);
    const className = React.isValidElement<{ className?: string }>(child)
        ? child.props.className || ""
        : "";

    return className.match(/language-([\w-]+)/)?.[1] || "text";
};

const copyText = async (text: string) => {
    if (navigator.clipboard?.writeText) {
        await navigator.clipboard.writeText(text);
        return;
    }

    const textarea = document.createElement("textarea");
    textarea.value = text;
    textarea.style.position = "fixed";
    textarea.style.left = "-9999px";
    textarea.style.top = "0";
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    document.execCommand("copy");
    document.body.removeChild(textarea);
};

const BlogCodeBlock: React.FC<BlogCodeBlockProps> = ({
    children,
    signatureText,
    ...props
}) => {
    const [copied, setCopied] = React.useState(false);
    const language = getCodeLanguage(children);
    const codeText = getTextContent(children).replace(/\n$/, "");

    const handleCopy = async () => {
        await copyText(codeText);
        setCopied(true);
        window.setTimeout(() => setCopied(false), 1600);
    };

    return (
        <div className="blog-code-block">
            <div className="blog-code-toolbar">
                <span className="blog-code-language">{language}</span>
                <button
                    type="button"
                    className="blog-code-copy"
                    onClick={handleCopy}
                    aria-label={copied ? "Copied code" : "Copy code"}
                    title={copied ? "Copied" : "Copy code"}
                >
                    {copied ? <CheckOutlined /> : <CopyOutlined />}
                    <span>{copied ? "Copied" : "Copy"}</span>
                </button>
            </div>
            <pre {...props}>{children}</pre>
            <span className="blog-code-signature">{signatureText}</span>
        </div>
    );
};

export const MarkdownPreview: React.FC<MarkdownPreviewProps> = ({
    content,
    className = "",
    headingLevelOffset = 0,
    theme = "default",
    signatureText,
}) => {
    const headingIdsByLine = React.useMemo(
        () => extractHeadingIdsByLine(content),
        [content]
    );
    const getHeadingId = (children: React.ReactNode, node: unknown) => {
        const line = (node as { position?: { start?: { line?: number } } })?.position?.start?.line;
        return (line ? headingIdsByLine.get(line) : undefined) || genHeadingId(String(children));
    };
    const isBlogTheme = theme === "blog";
    const resolvedClassName = [
        "markdown-body",
        isBlogTheme ? "blog-article-theme" : "",
        className,
    ].filter(Boolean).join(" ");
    const displaySignature = signatureText?.trim() || "FNote";
    const Heading = ({
        level,
        children,
        node,
        ...props
    }: {
        level: number;
        children: React.ReactNode;
        node: unknown;
    } & React.HTMLAttributes<HTMLHeadingElement>) => {
        const resolvedLevel = Math.min(6, level + Math.max(0, headingLevelOffset));
        const Tag = `h${resolvedLevel}` as "h1" | "h2" | "h3" | "h4" | "h5" | "h6";
        return <Tag id={getHeadingId(children, node)} {...props}>{children}</Tag>;
    };
    const components: Components = {
        h1: ({children, node, ...props}) => <Heading level={1} node={node} {...props}>{children}</Heading>,
        h2: ({children, node, ...props}) => <Heading level={2} node={node} {...props}>{children}</Heading>,
        h3: ({children, node, ...props}) => <Heading level={3} node={node} {...props}>{children}</Heading>,
        h4: ({children, node, ...props}) => <Heading level={4} node={node} {...props}>{children}</Heading>,
        h5: ({children, node, ...props}) => <Heading level={5} node={node} {...props}>{children}</Heading>,
        h6: ({children, node, ...props}) => <Heading level={6} node={node} {...props}>{children}</Heading>,
        ...(isBlogTheme
            ? {
                blockquote: ({children, ...props}) => (
                    <blockquote {...props}>
                        <div className="blog-quote-signature">{displaySignature}</div>
                        {children}
                    </blockquote>
                ),
                pre: ({children, node: _node, ...props}) => (
                    <BlogCodeBlock {...props} signatureText={displaySignature}>
                        {children}
                    </BlogCodeBlock>
                ),
            }
            : {}),
    };

    return (
        <div className={resolvedClassName}>
            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeRaw, rehypeHighlight]}
                components={components}
            >
                {content}
            </ReactMarkdown>
        </div>
    );
};
