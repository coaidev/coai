import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog.tsx";
import { Image, MenuSquare, PanelRight, ClipboardPaste } from "lucide-react";
import { useTranslation } from "react-i18next";
import "@/assets/common/editor.less";
import MarkdownExport from "./Markdown.tsx";
import React, { useMemo, useState } from "react";
import { Toggle } from "./ui/toggle.tsx";
import { mobile } from "@/utils/device.ts";
import { ChatAction } from "@/components/home/assemblies/ChatAction.tsx";
import { cn } from "@/components/ui/lib/utils.ts";
import { Message } from "@/api/types.tsx";
import {
  useMessages,
} from "@/store/chat.ts";

type ExportAllMsgButtonProps = {
  maxLength?: number;

  formatter?: (value: string) => string;
  title?: string;

  open?: boolean;
  setOpen?: (open: boolean) => void;
  children?: React.ReactNode;

  submittable?: boolean;
  onSubmit?: (value: string) => void;
  closeOnSubmit?: boolean;
};

function ExportAllMsgButtonCall({
  formatter,
}: ExportAllMsgButtonProps) {
  const [openPreview, setOpenPreview] = useState(!mobile);
  const [openInput, setOpenInput] = useState(true);
  const messages: Message[] = useMessages();
  const jsonArray = messages.map(message => ({
    role: message.role,
    content: message.content
  }));
  const jsonString = `\`\`\`json\n${JSON.stringify(jsonArray, null, 4)}\n\`\`\``;

  function convertToMarkdown(jsonArray: { role: any; content: any; }[]) {
    const { t } = useTranslation();
    return jsonArray.map(({ role, content }) => {
      const roleText = role === 'user' ? t("export.user_says") : t("export.ai_says");
      return `## ${roleText}\n\n${content}`;
    }).join('\n\n');
  }

  const markdownString = convertToMarkdown(jsonArray);
  const markdownValue = useMemo(() => {
    return formatter ? formatter(markdownString) : markdownString;
  }, [markdownString, formatter]);




  return (
    <div className={`editor-container`}>
      <div className={`editor-toolbar`}>
        <div className={`grow`} />
        <Toggle
          variant={`outline`}
          className={`h-8 w-8 p-0`}
          pressed={openInput && !openPreview}
          onClick={() => {
            setOpenPreview(false);
            setOpenInput(true);
          }}
        >
          <MenuSquare className={`h-3.5 w-3.5`} />
        </Toggle>

        <Toggle
          variant={`outline`}
          className={`h-8 w-8 p-0`}
          pressed={openInput && openPreview}
          onClick={() => {
            setOpenPreview(true);
            setOpenInput(true);
          }}
        >
          <PanelRight className={`h-3.5 w-3.5`} />
        </Toggle>

        <Toggle
          variant={`outline`}
          className={`h-8 w-8 p-0`}
          pressed={!openInput && openPreview}
          onClick={() => {
            setOpenPreview(true);
            setOpenInput(false);
          }}
        >
          <Image className={`h-3.5 w-3.5`} />
        </Toggle>
      </div>
      <div className={`editor-wrapper`}>
        <div
          className={cn(
            "editor-object",
            openInput && "show-editor",
            openPreview && "show-preview",
          )}
        >
          {openInput && (
            <MarkdownExport
              className={cn(
                `transition-all`
              )}
              loading={true}
              children={jsonString}

            />


          )}
          {openPreview && (
            <MarkdownExport
              loading={true}
              children={markdownValue}
            />
          )}
        </div>
      </div>

    </div>
  );
}

function ExportAllMsgButton(props: ExportAllMsgButtonProps) {
  const { t } = useTranslation();

  return (
    <>
      <Dialog open={props.open} onOpenChange={props.setOpen}>
        {!props.setOpen && (
          <DialogTrigger asChild>
            {props.children ?? (
              <ChatAction text={t("export.text")}>
                <ClipboardPaste className={`h-4 w-4`} />
              </ChatAction>
            )}
          </DialogTrigger>
        )}
        <DialogContent className={`editor-dialog flex-dialog`}>
          <DialogHeader>
            <DialogTitle>{props.title ?? t("export.user_used")}</DialogTitle>
            <DialogDescription asChild>
              <ExportAllMsgButtonCall {...props}>
              </ExportAllMsgButtonCall>
            </DialogDescription>
          </DialogHeader>
        </DialogContent>
      </Dialog>
    </>
  );
}

export default ExportAllMsgButton;

export function JSONTransMarkdownProvider({ ...props }: ExportAllMsgButtonProps) {
  return (
    <ExportAllMsgButton
      {...props}
      formatter={(value) => `\`\`\`markdown\n${value}\n\`\`\``}

    />
  );
}
