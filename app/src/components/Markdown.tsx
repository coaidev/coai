import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkBreaks from "remark-breaks";
import rehypeKatex from "rehype-katex";
import rehypeRaw from "rehype-raw";
import "@/assets/markdown/all.less";
import { useEffect, useMemo } from "react";
import { cn } from "@/components/ui/lib/utils.ts";
import Label from "@/components/markdown/Label.tsx";
import Link from "@/components/markdown/Link.tsx";
import Code, { CodeProps } from "@/components/markdown/Code.tsx";

type MarkdownProps = {
  children: string;
  className?: string;
  acceptHtml?: boolean;
  codeStyle?: string;
  loading?: boolean;
};

function MarkdownContent({
  children,
  className,
  acceptHtml,
  codeStyle,
  loading,
}: MarkdownProps) {
  useEffect(() => {
    document.querySelectorAll(".file-instance").forEach((el) => {
      const parent = el.parentElement as HTMLElement;
      if (!parent.classList.contains("file-block"))
        parent.classList.add("file-block");
    });
  }, [children]);

  const rehypePlugins = useMemo(() => {
    const plugins = [rehypeKatex as any];
    return acceptHtml ? [...plugins, rehypeRaw] : plugins;
  }, [acceptHtml]);

  const components = useMemo(() => {
    return {
      p: Label,
      a: Link,
      code: (props: CodeProps) => (
        <Code {...props} loading={loading} codeStyle={codeStyle} />
      ),
    };
  }, [codeStyle]);

  return (
    <ReactMarkdown
      remarkPlugins={[remarkMath, remarkGfm, remarkBreaks]}
      rehypePlugins={rehypePlugins}
      className={cn("markdown-body", className)}
      children={children}
      skipHtml={acceptHtml}
      components={components}
    />
  );
}

function Markdown({
  children,
  acceptHtml,
  codeStyle,
  className,
  loading,
}: MarkdownProps) {
  const processedContent = useMemo(() => {
    let content = children;
    
    // Inline math: replace \(...\) with $ ... $
    content = content.replace(/\\\((.*?)\\\)/g, (_, equation) => `$ ${equation.trim()} $`);
    
    // Block math: replace \[...\] with $$...$$ on separate lines
    content = content.replace(
      /\s*\\\[\s*([\s\S]*?)\s*\\\]\s*/g,
      (_, equation) => `\n$$\n${equation.trim()}\n$$\n`
    );
    
    return content;
  }, [children]);

  return useMemo(
    () => (
      <MarkdownContent
        children={processedContent}
        acceptHtml={acceptHtml}
        codeStyle={codeStyle}
        className={className}
        loading={loading}
      />
    ),
    [processedContent, acceptHtml, codeStyle, className, loading],
  );
}
type CodeMarkdownProps = MarkdownProps & {
  filename: string;
};

export function CodeMarkdown({ filename, ...props }: CodeMarkdownProps) {
  const suffix = filename.includes(".") ? filename.split(".").pop() : "";
  const children = useMemo(() => {
    const content = props.children.toString();

    return `\`\`\`${suffix}\n${content}\n\`\`\``;
  }, [props.children]);

  return <Markdown {...props}>{children}</Markdown>;
}

export default Markdown;
