import { useTranslation } from "react-i18next";
import { useEffect, useMemo, useReducer, useRef, useState } from "react";
import FileAction from "@/components/FileProvider.tsx";
import { useSelector } from "react-redux";
import { selectAuthenticated, selectInit } from "@/store/auth.ts";
import {
  listenMessageEvent,
  selectCurrent,
  selectModel,
  selectSupportModels,
  useMessageActions,
  useMessages,
  useWorking,
} from "@/store/chat.ts";
import { formatMessage } from "@/utils/processor.ts";
import ChatInterface from "@/components/home/ChatInterface.tsx";
import EditorAction from "@/components/EditorProvider.tsx";
import ModelFinder from "./ModelFinder.tsx";
import { clearHistoryState, getQueryParam } from "@/utils/path.ts";
import { forgetMemory, popMemory } from "@/utils/memory.ts";
import { useToast } from "@/components/ui/use-toast.ts";
import { ToastAction } from "@/components/ui/toast.tsx";
import { alignSelector } from "@/store/settings.ts";
import { FileArray } from "@/api/file.ts";
import {
  MarketAction,
  MaskAction,
  SettingsAction,
  WebAction,
} from "@/components/home/assemblies/ChatAction.tsx";
import ChatSpace from "@/components/home/ChatSpace.tsx";
import ActionButton from "@/components/home/assemblies/ActionButton.tsx";
import ChatInput from "@/components/home/assemblies/ChatInput.tsx";
import ScrollAction from "@/components/home/assemblies/ScrollAction.tsx";
import { cn } from "@/components/ui/lib/utils.ts";
import { goAuth } from "@/utils/app.ts";
import { getModelFromId } from "@/conf/model.ts";
import { JSONTransMarkdownProvider } from "@/components/ExportAllMsgButton.tsx";

type InterfaceProps = {
  scrollable: boolean;
  setTarget: (instance: HTMLElement | null) => void;
};

function Interface(props: InterfaceProps) {
  const messages = useMessages();
  return messages.length > 0 ? <ChatInterface {...props} /> : <ChatSpace />;
}

function fileReducer(state: FileArray, action: Record<string, any>): FileArray {
  switch (action.type) {
    case "add":
      return [...state, action.payload];
    case "remove":
      return state.filter((_, i) => i !== action.payload);
    case "clear":
      return [];
    default:
      return state;
  }
}

function ChatWrapper() {
  const { t } = useTranslation();
  const { toast } = useToast();
  const { send: sendAction } = useMessageActions();
  const process = listenMessageEvent();
  const [files, fileDispatch] = useReducer(fileReducer, []);
  const [input, setInput] = useState("");
  const [visible, setVisibility] = useState(false);
  const init = useSelector(selectInit);
  const current = useSelector(selectCurrent);
  const auth = useSelector(selectAuthenticated);
  const model = useSelector(selectModel);
  const target = useRef(null);
  const align = useSelector(alignSelector);

  const working = useWorking();
  const supportModels = useSelector(selectSupportModels);

  const requireAuth = useMemo(
    (): boolean => !!getModelFromId(supportModels, model)?.auth,
    [model],
  );

  const [instance, setInstance] = useState<HTMLElement | null>(null);

  function clearFile() {
    fileDispatch({ type: "clear" });
  }

  async function processSend(
    data: string,
    passAuth?: boolean,
  ): Promise<boolean> {
    if (requireAuth && !auth && !passAuth) {
      toast({
        title: t("login-require"),
        action: (
          <ToastAction altText={t("login")} onClick={goAuth}>
            {t("login")}
          </ToastAction>
        ),
      });
      return false;
    }

    if (working) return false;

    const message: string = formatMessage(files, data);
    if (message.length > 0 && data.trim().length > 0) {
      if (await sendAction(message)) {
        forgetMemory("history");
        clearFile();
        return true;
      }
    }
    return false;
  }

  async function handleSend() {
    // because of the function wrapper, we need to update the selector state using props.
    if (await processSend(input)) {
      setInput("");
    }
  }

  async function handleCancel() {
    process({ id: current, event: "stop" });
  }

  useEffect(() => {
    window.addEventListener("load", () => {
      const el = document.getElementById("input");
      if (el) el.focus();
    });
  }, []);

  useEffect(() => {
    if (!init) return;
    const query = getQueryParam("q").trim();
    if (query.length > 0) processSend(query).then();
    clearHistoryState();
  }, [init]);

  useEffect(() => {
    const history: string = popMemory("history");
    if (history.length) {
      setInput(history);
      toast({
        title: t("chat.recall"),
        description: t("chat.recall-desc"),
        action: (
          <ToastAction
            altText={t("chat.recall-cancel")}
            onClick={() => {
              setInput("");
            }}
          >
            {t("chat.recall-cancel")}
          </ToastAction>
        ),
      });
    }
  }, []);

  return (
    <div className={`chat-container`}>
      <div className={`chat-wrapper`}>
        <Interface setTarget={setInstance} scrollable={!visible} />
        <div className={`chat-input`}>
          <div className={`input-action`}>
            <ScrollAction
              visible={visible}
              setVisibility={setVisibility}
              target={instance}
            />
            <WebAction visible={!visible} />
            <FileAction files={files} dispatch={fileDispatch} />
            <EditorAction value={input} onChange={setInput} />
            <MaskAction />
            <MarketAction />
            <SettingsAction />
            <JSONTransMarkdownProvider />
          </div>
          <div className={`input-wrapper`}>
            <div className={`chat-box no-scrollbar`}>
              <ChatInput
                className={cn(align && "align")}
                target={target}
                value={input}
                onValueChange={setInput}
                onEnterPressed={handleSend}
              />
            </div>
            <ActionButton
              working={working}
              onClick={() => (working ? handleCancel() : handleSend())}
            />
          </div>
          <div className={`input-options`}>
            <ModelFinder side={`bottom`} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default ChatWrapper;
