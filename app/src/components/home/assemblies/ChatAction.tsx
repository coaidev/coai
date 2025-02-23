import {
  openMarket,
  openMask,
  selectModel,
  selectSupportModels,
  selectWeb,
  toggleWeb,
  useConversationActions,
} from "@/store/chat.ts";
import {
  Blocks,
  Globe,
  LandPlot,
  Play,
  Search,
  Settings,
  Wand2,
} from "lucide-react";
import { useDispatch, useSelector } from "react-redux";
import { useTranslation } from "react-i18next";
import React, { useEffect, useMemo, useRef } from "react";
import { openDialog } from "@/store/settings.ts";
import { cn } from "@/components/ui/lib/utils.ts";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { ScrollArea } from "@/components/ui/scroll-area.tsx";
import { getModelAvatar } from "@/components/home/ModelMarket.tsx";
import { Button } from "@/components/ui/button.tsx";
import Tips from "@/components/Tips.tsx";
import { Model } from "@/api/types.tsx";
import { Input } from "@/components/ui/input.tsx";
import { splitList } from "@/utils/base.ts";

type ChatActionProps = {
  className?: string;
  text?: string;
  children?: React.ReactNode;
  onClick?: () => void;
};
export const ChatAction = React.forwardRef<HTMLDivElement, ChatActionProps>(
  (props, ref) => {
    const { className, text, children, onClick } = props;
    const label = useRef(null);
    const [labelWidth, setLabelWidth] = React.useState(0);

    useEffect(() => {
      if (!label.current) return;

      const target = label.current as HTMLDivElement;
      const width = target.clientWidth || target.offsetWidth;
      setLabelWidth(width);
    }, [text, label]);

    return (
      <div
        className={cn("action chat-action", className)}
        onClick={onClick}
        ref={ref}
        style={{ "--width": `${labelWidth}px` } as React.CSSProperties}
      >
        {children}
        <div className="text" ref={label}>
          {text}
        </div>
      </div>
    );
  },
);

type WebActionProps = {
  visible: boolean;
};

export function WebAction({ visible }: WebActionProps) {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const web = useSelector(selectWeb);

  return (
    visible && (
      <ChatAction
        className={cn(web && "active")}
        text={t("chat.web")}
        onClick={() => dispatch(toggleWeb())}
      >
        <Globe className={cn("h-4 w-4 web", web && "enable")} />
      </ChatAction>
    )
  );
}

export function SettingsAction() {
  const { t } = useTranslation();
  const dispatch = useDispatch();

  return (
    <ChatAction
      text={t("settings.description")}
      onClick={() => dispatch(openDialog())}
    >
      <Settings className={`h-4 w-4`} />
    </ChatAction>
  );
}
export function ExportButtonAction() {
  const { t } = useTranslation();
  const dispatch = useDispatch();

  return (
    <ChatAction
      text={t("export.markdown")}
      onClick={() => dispatch(openDialog())}
    >
      <Settings className={`h-4 w-4`} />
    </ChatAction>
  );
}

export function MarketAction() {
  const { t } = useTranslation();
  const dispatch = useDispatch();

  const current = useSelector(selectModel);
  const supportModels = useSelector(selectSupportModels);

  const [search, setSearch] = React.useState("");
  const [open, setOpen] = React.useState(false);

  const { selected } = useConversationActions();

  const models = useMemo(() => {
    const content = search.trim();
    if (content === "") return supportModels;

    const raw = splitList(search.toLowerCase(), [" ", ",", ";", "-"]);
    return supportModels.filter((model) => {
      return raw.every((item) => {
        return model.name.toLowerCase().includes(item);
      });
    });
  }, [supportModels, search]);

  const event = (model: Model) => {
    selected(model.id);
    setOpen(false);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <ChatAction text={t("market.title")}>
          <Blocks className={`h-4 w-4`} />
        </ChatAction>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("market.list")}</DialogTitle>
          <DialogDescription asChild>
            <div>
              <Button
                size={`default-sm`}
                className={`mt-2 w-full`}
                onClick={() => dispatch(openMarket())}
              >
                <LandPlot className={`h-4 w-4 mr-2`} />
                {t("market.go")}
              </Button>
              <div className={`flex flex-row mt-4 mb-2`}>
                <Input
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  placeholder={t("market.search")}
                />
                <Button
                  className={`ml-2 shrink-0`}
                  variant={`outline`}
                  size={`icon`}
                >
                  <Search className={`h-4 w-4`} />
                </Button>
              </div>
              <ScrollArea
                className={`flex flex-col mt-2.5 border rounded-md max-h-[60vh]`}
              >
                {models.length > 0 ? (
                  models.map((model, index) => (
                    <div
                      key={index}
                      onClick={() => event(model)}
                      className={cn(
                        "flex sm:flex-row flex-col items-center px-4 py-4 sm:py-2 border-b last:border-none select-none cursor-pointer transition-all hover:bg-background-container",
                        model.id === current && "bg-background-container",
                      )}
                    >
                      <img
                        className={`w-6 h-6 border rounded-md`}
                        src={getModelAvatar(model.avatar)}
                        alt={""}
                      />
                      <div className={`text-common sm:ml-2 mt-2 sm:mt-0`}>
                        {model.name}
                      </div>
                      <Tips className={`hidden sm:inline-block`}>
                        {model.id}
                      </Tips>
                      <Button
                        className={cn(
                          "ml-auto",
                          model.id === current && "text-common",
                        )}
                        size={`icon-sm`}
                        variant={`ghost`}
                        onClick={() => event(model)}
                      >
                        <Play className={`h-3.5 w-3.5`} />
                      </Button>
                    </div>
                  ))
                ) : (
                  <div
                    className={`flex justify-center items-center py-8 select-none`}
                  >
                    {t("empty")}
                  </div>
                )}
              </ScrollArea>
            </div>
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  );
}

export function MaskAction() {
  const { t } = useTranslation();
  const dispatch = useDispatch();

  return (
    <ChatAction text={t("mask.title")} onClick={() => dispatch(openMask())}>
      <Wand2 className={`h-4 w-4`} />
    </ChatAction>
  );
}
