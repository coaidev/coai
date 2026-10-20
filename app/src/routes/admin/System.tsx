import { useTranslation } from "react-i18next";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import Paragraph, {
  ParagraphDescription,
  ParagraphFooter,
  ParagraphItem,
  ParagraphSpace,
} from "@/components/Paragraph.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx";
import { Label } from "@/components/ui/label.tsx";
import { Input } from "@/components/ui/input.tsx";
import { useMemo, useReducer, useState } from "react";
import { formReducer } from "@/utils/form.ts";
import { NumberInput } from "@/components/ui/number-input.tsx";
import {
  CommonState,
  commonWhiteList,
  GeneralState,
  getConfig,
  initialSystemState,
  MailState,
  SearchState,
  setConfig,
  SiteState,
  SystemProps,
  testWebSearching,
  updateRootPassword,
} from "@/admin/api/system.ts";
import { useEffectAsync } from "@/utils/hook.ts";
import { withNotify } from "@/api/common.ts";
import { doVerify } from "@/api/auth.ts";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { DialogTitle } from "@radix-ui/react-dialog";
import Require from "@/components/Require.tsx";
import { Loader2, PencilLine, RotateCw, Save, Settings2 } from "lucide-react";
import { FlexibleTextarea, Textarea } from "@/components/ui/textarea.tsx";
import Tips from "@/components/Tips.tsx";
import { cn } from "@/components/ui/lib/utils.ts";
import { Switch } from "@/components/ui/switch.tsx";
import { MultiCombobox } from "@/components/ui/multi-combobox.tsx";
import { allGroups } from "@/utils/groups.ts";
import { useChannelModels } from "@/admin/hook.tsx";
import { useSelector } from "react-redux";
import { selectSupportModels } from "@/store/chat.ts";
import { JSONEditorProvider } from "@/components/EditorProvider.tsx";
import { Combobox } from "@/components/ui/combo-box.tsx";

type CompProps<T> = {
  data: T;
  form: SystemProps;
  dispatch: (action: any) => void;
  onChange: (doToast?: boolean) => Promise<void>;
};

function RootDialog() {
  const { t } = useTranslation();
  const [open, setOpen] = useState<boolean>(false);
  const [password, setPassword] = useState<string>("");
  const [repeat, setRepeat] = useState<string>("");

  const onPost = async () => {
    const res = await updateRootPassword(password);
    withNotify(t, res, true);
    if (res.status) {
      setPassword("");
      setRepeat("");
      setOpen(false);

      setTimeout(() => {
        window.location.reload();
      }, 1000);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant={`outline`} size={`sm`}>
          {t("admin.system.updateRoot")}
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("admin.system.updateRoot")}</DialogTitle>
          <DialogDescription>
            <div className={`mb-4 select-none`}>
              {t("admin.system.updateRootTip")}
            </div>
            <Input
              className={`mb-2`}
              type={`password`}
              placeholder={t("admin.system.updateRootPlaceholder")}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Input
              type={`password`}
              placeholder={t("admin.system.updateRootRepeatPlaceholder")}
              value={repeat}
              onChange={(e) => setRepeat(e.target.value)}
            />
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant={`outline`}
            onClick={() => {
              setPassword("");
              setRepeat("");
              setOpen(false);
            }}
          >
            {t("admin.cancel")}
          </Button>
          <Button
            variant={`default`}
            loading={true}
            onClick={onPost}
            disabled={
              password.trim().length === 0 || password.trim() !== repeat.trim()
            }
          >
            {t("admin.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function General({ data, dispatch, onChange }: CompProps<GeneralState>) {
  const { t } = useTranslation();

  return (
    <Paragraph
      title={t("admin.system.general")}
      configParagraph={true}
      isCollapsed={true}
    >
      <ParagraphItem>
        <Label>{t("admin.system.title")}</Label>
        <Input
          value={data.title}
          onChange={(e) =>
            dispatch({
              type: "update:general.title",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.titleTip")}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.docs")}</Label>
        <Input
          value={data.docs}
          onChange={(e) =>
            dispatch({
              type: "update:general.docs",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.docsTip")}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.logo")}</Label>
        <Input
          value={data.logo}
          onChange={(e) =>
            dispatch({
              type: "update:general.logo",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.logoTip", {
            logo: `${window.location.protocol}//${window.location.host}/favicon.ico`,
          })}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.backend")}</Label>
        <Input
          value={data.backend}
          onChange={(e) =>
            dispatch({
              type: "update:general.backend",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.backendPlaceholder")}
        />
      </ParagraphItem>
      <ParagraphDescription border>
        {t("admin.system.backendTip", {
          backend: `${window.location.protocol}//${window.location.host}/api`,
        })}
      </ParagraphDescription>
      <ParagraphItem>
        <Label>{t("admin.system.file")}</Label>
        <Input
          value={data.file}
          onChange={(e) =>
            dispatch({
              type: "update:general.file",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.filePlaceholder")}
        />
      </ParagraphItem>
      <ParagraphDescription border>
        {t("admin.system.fileTip")}
      </ParagraphDescription>
      <ParagraphItem>
        <Label>PWA Manifest</Label>
        <JSONEditorProvider
          value={data.pwa_manifest ?? ""}
          onChange={(value) =>
            dispatch({ type: "update:general.pwa_manifest", value })
          }
        >
          <Button variant={`outline`}>
            <PencilLine className={`h-4 w-4 mr-1`} />
            {t("edit")}
          </Button>
        </JSONEditorProvider>
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          {t("admin.system.debugMode")}
          <Tips
            className={`inline-block`}
            content={t("admin.system.debugModeTip")}
          />
        </Label>
        <Switch
          checked={data.debug_mode}
          onCheckedChange={(value) => {
            dispatch({ type: "update:general.debug_mode", value });
          }}
        />
      </ParagraphItem>
      <ParagraphSpace />
      <ParagraphFooter>
        <div className={`grow`} />
        <RootDialog />
        <Button
          size={`sm`}
          loading={true}
          onClick={async () => await onChange()}
        >
          {t("admin.system.save")}
        </Button>
      </ParagraphFooter>
    </Paragraph>
  );
}

function Mail({ data, dispatch, onChange }: CompProps<MailState>) {
  const { t } = useTranslation();
  const [email, setEmail] = useState<string>("");

  const [mailDialog, setMailDialog] = useState<boolean>(false);

  const valid = useMemo((): boolean => {
    return (
      data.host.length > 0 &&
      data.port > 0 &&
      data.port < 65535 &&
      data.username.length > 0 &&
      data.password.length > 0 &&
      data.from.length > 0
    );
  }, [data]);

  const onTest = async () => {
    if (!email.trim()) return;
    await onChange(false);
    const res = await doVerify(email);
    withNotify(t, res, true);

    if (res.status) setMailDialog(false);
  };

  const white_list = useMemo(() => {
    const raw = data.white_list.custom
      .split(",")
      .map((item) => item.trim())
      .filter((item) => item.length > 0);

    return [...commonWhiteList, ...raw];
  }, [data]);

  return (
    <Paragraph
      title={t("admin.system.mail")}
      configParagraph={true}
      isCollapsed={true}
    >
      {!valid && (
        <ParagraphDescription border={true}>
          {t("admin.system.mailConfNotValid")}
        </ParagraphDescription>
      )}
      <ParagraphItem>
        <Label>
          <Require /> {t("admin.system.mailHost")}
        </Label>
        <Input
          value={data.host}
          onChange={(e) =>
            dispatch({
              type: "update:mail.host",
              value: e.target.value,
            })
          }
          placeholder={`smtp.qcloudmail.com`}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          <Require /> {t("admin.system.mailProtocol")}
        </Label>
        <Select
          value={data.protocol ? "true" : "false"}
          onValueChange={(value: string) => {
            dispatch({
              type: "update:mail.protocol",
              value: value === "true",
            });
          }}
        >
          <SelectTrigger className={`select`}>
            <SelectValue
              placeholder={
                data.protocol
                  ? t("admin.system.mailProtocolTLS")
                  : t("admin.system.mailProtocolSSL")
              }
            />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="true">TLS</SelectItem>
            <SelectItem value="false">SSL</SelectItem>
          </SelectContent>
        </Select>
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          <Require /> {t("admin.system.mailPort")}
        </Label>
        <NumberInput
          value={data.port}
          onValueChange={(value) =>
            dispatch({ type: "update:mail.port", value })
          }
          placeholder={`465`}
          min={0}
          max={65535}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          <Require /> {t("admin.system.mailUser")}
        </Label>
        <Input
          value={data.username}
          onChange={(e) =>
            dispatch({
              type: "update:mail.username",
              value: e.target.value,
            })
          }
          className={cn("transition-all duration-300")}
          placeholder={t("admin.system.mailUser")}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          <Require /> {t("admin.system.mailPass")}
        </Label>
        <Input
          value={data.password}
          onChange={(e) =>
            dispatch({
              type: "update:mail.password",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.mailPass")}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          <Require /> {t("admin.system.mailFrom")}
        </Label>
        <Input
          value={data.from}
          onChange={(e) =>
            dispatch({
              type: "update:mail.from",
              value: e.target.value,
            })
          }
          placeholder={`${t("admin.system.mailFrom")} <${data.username}@${location.hostname}>`}
          className={cn("transition-all duration-300")}
        />
      </ParagraphItem>
      <ParagraphSpace />
      <ParagraphItem>
        <Label>{t("admin.system.mailEnableWhitelist")}</Label>
        <Switch
          checked={data.white_list.enabled}
          onCheckedChange={(value) => {
            dispatch({
              type: "update:mail.white_list.enabled",
              value,
            });
          }}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.mailWhitelist")}</Label>
        <MultiCombobox
          value={data.white_list.white_list}
          list={white_list}
          disabled={!data.white_list.enabled}
          onChange={(value) => {
            dispatch({
              type: "update:mail.white_list.white_list",
              value,
            });
          }}
          placeholder={t("admin.system.mailWhitelistSelected", {
            length: data.white_list.white_list.length,
          })}
          searchPlaceholder={t("admin.system.mailWhitelistSearchPlaceholder")}
        />
      </ParagraphItem>
      <Input
        className={`mb-2`}
        value={data.white_list.custom}
        onChange={(e) =>
          dispatch({
            type: "update:mail.white_list.custom",
            value: e.target.value,
          })
        }
        disabled={!data.white_list.enabled}
        placeholder={t("admin.system.customWhitelistPlaceholder")}
      />
      <ParagraphFooter>
        <div className={`grow`} />
        <Dialog open={mailDialog} onOpenChange={setMailDialog}>
          <DialogTrigger asChild>
            <Button variant={`outline`} size={`sm`}>
              {t("admin.system.test")}
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("admin.system.test")}</DialogTitle>
              <DialogDescription className={`pt-2`}>
                <Input
                  placeholder={t("auth.email-placeholder")}
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button
                variant={`outline`}
                onClick={() => {
                  setEmail("");
                  setMailDialog(false);
                }}
              >
                {t("admin.cancel")}
              </Button>
              <Button variant={`default`} loading={true} onClick={onTest}>
                {t("admin.confirm")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        <Button
          size={`sm`}
          loading={true}
          onClick={async () => await onChange()}
        >
          {t("admin.system.save")}
        </Button>
      </ParagraphFooter>
    </Paragraph>
  );
}

function Site({ data, dispatch, onChange }: CompProps<SiteState>) {
  const { t } = useTranslation();

  return (
    <Paragraph
      title={t("admin.system.site")}
      configParagraph={true}
      isCollapsed={true}
    >
      <ParagraphItem>
        <Label>
          {t("admin.system.closeRegistration")}
          <Tips
            className={`inline-block`}
            content={t("admin.system.closeRegistrationTip")}
          />
        </Label>
        <Switch
          checked={data.close_register}
          onCheckedChange={(value) => {
            dispatch({ type: "update:site.close_register", value });
          }}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          {t("admin.system.closeRelay")}
          <Tips
            className={`inline-block`}
            content={t("admin.system.closeRelayTip")}
          />
        </Label>
        <Switch
          checked={data.close_relay}
          onCheckedChange={(value) => {
            dispatch({ type: "update:site.close_relay", value });
          }}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          {t("admin.system.relayPlan")}
          <Tips
            className={`inline-block`}
            content={t("admin.system.relayPlanTip")}
          />
        </Label>
        <Switch
          checked={data.relay_plan}
          onCheckedChange={(value) => {
            dispatch({ type: "update:site.relay_plan", value });
          }}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.quota")}
          <Tips content={t("admin.system.quotaTip")} />
        </Label>
        <NumberInput
          value={data.quota}
          onValueChange={(value) =>
            dispatch({ type: "update:site.quota", value })
          }
          placeholder={`5`}
          min={0}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.buyLink")}</Label>
        <Input
          value={data.buy_link}
          onChange={(e) =>
            dispatch({
              type: "update:site.buy_link",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.buyLinkPlaceholder")}
        />
      </ParagraphItem>
      <ParagraphItem rowLayout={true}>
        <Label>{t("admin.system.announcement")}</Label>
        <FlexibleTextarea
          value={data.announcement}
          rows={12}
          onChange={(e) =>
            dispatch({
              type: "update:site.announcement",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.announcementPlaceholder")}
        />
      </ParagraphItem>
      <ParagraphItem rowLayout={true}>
        <Label>{t("admin.system.contact")}</Label>
        <FlexibleTextarea
          value={data.contact}
          rows={6}
          onChange={(e) =>
            dispatch({
              type: "update:site.contact",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.contactPlaceholder")}
        />
      </ParagraphItem>
      <ParagraphSpace />
      <ParagraphItem rowLayout={true}>
        <Label>{t("admin.system.footer")}</Label>
        <FlexibleTextarea
          rows={6}
          value={data.footer}
          onChange={(e) =>
            dispatch({
              type: "update:site.footer",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.footerPlaceholder")}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.authFooter")}</Label>
        <Switch
          checked={data.auth_footer}
          onCheckedChange={(value) => {
            dispatch({ type: "update:site.auth_footer", value });
          }}
        />
      </ParagraphItem>
      <ParagraphFooter>
        <div className={`grow`} />
        <Button
          size={`sm`}
          loading={true}
          onClick={async () => await onChange()}
        >
          {t("admin.system.save")}
        </Button>
      </ParagraphFooter>
    </Paragraph>
  );
}

function Common({ form, data, dispatch, onChange }: CompProps<CommonState>) {
  const { t } = useTranslation();

  const { channelModels } = useChannelModels();
  const supportModels = useSelector(selectSupportModels);

  return (
    <Paragraph
      title={t("admin.system.common")}
      configParagraph={true}
      isCollapsed={true}
    >
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.image_store")}
          <Tips content={t("admin.system.image_storeTip")} />
        </Label>
        <Switch
          checked={data.image_store}
          onCheckedChange={(value) => {
            dispatch({ type: "update:common.image_store", value });
          }}
        />
      </ParagraphItem>
      {data.image_store && form.general.backend.length === 0 && (
        <ParagraphDescription border={true}>
          {t("admin.system.image_storeNoBackend")}
        </ParagraphDescription>
      )}
      <ParagraphSpace />
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.cache")}
          <Tips content={t("admin.system.cacheTip")} />
        </Label>
        <MultiCombobox
          value={data.cache}
          onChange={(value) => {
            dispatch({ type: "update:common.cache", value });
          }}
          list={channelModels}
          placeholder={t("admin.system.cachePlaceholder", {
            length: (data.cache ?? []).length,
          })}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          {t("admin.system.cacheExpired")}
          <Tips
            className={`inline-block`}
            content={t("admin.system.cacheExpiredTip")}
          />
        </Label>
        <NumberInput
          value={data.expire}
          onValueChange={(value) =>
            dispatch({ type: "update:common.expire", value })
          }
          min={0}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>
          {t("admin.system.cacheSize")}
          <Tips
            className={`inline-block`}
            content={t("admin.system.cacheSizeTip")}
          />
        </Label>
        <NumberInput
          value={data.size}
          onValueChange={(value) =>
            dispatch({ type: "update:common.size", value })
          }
          min={0}
        />
      </ParagraphItem>
      <ParagraphItem>
        <div className={`flex flex-row flex-wrap gap-2 ml-auto`}>
          <Button
            variant={`outline`}
            onClick={() => dispatch({ type: "update:common.cache", value: [] })}
          >
            <Settings2
              className={`inline-flex h-4 w-4 mr-2 translate-y-[1px]`}
            />
            {t("admin.system.cacheNone")}
          </Button>
          <Button
            variant={`outline`}
            onClick={() =>
              dispatch({
                type: "update:common.cache",
                value: supportModels
                  .filter((item) => item.free)
                  .map((item) => item.id),
              })
            }
          >
            <Settings2
              className={`inline-flex h-4 w-4 mr-2 translate-y-[1px]`}
            />
            {t("admin.system.cacheFree")}
          </Button>
          <Button
            variant={`outline`}
            onClick={() =>
              dispatch({ type: "update:common.cache", value: channelModels })
            }
          >
            <Settings2
              className={`inline-flex h-4 w-4 mr-2 translate-y-[1px]`}
            />
            {t("admin.system.cacheAll")}
          </Button>
        </div>
      </ParagraphItem>
      <ParagraphSpace />
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.article")}
          <Tips content={t("admin.system.articleTip")} />
        </Label>
        <MultiCombobox
          value={data.article}
          onChange={(value) => {
            dispatch({ type: "update:common.article", value });
          }}
          list={allGroups}
          listTranslate={`admin.channels.groups`}
          placeholder={t("admin.system.groupPlaceholder", {
            length: (data.article ?? []).length,
          })}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.generate")}
          <Tips content={t("admin.system.generateTip")} />
        </Label>
        <MultiCombobox
          value={data.generation}
          onChange={(value) => {
            dispatch({ type: "update:common.generation", value });
          }}
          list={allGroups}
          listTranslate={`admin.channels.groups`}
          placeholder={t("admin.system.groupPlaceholder", {
            length: (data.generation ?? []).length,
          })}
        />
      </ParagraphItem>
      <ParagraphFooter>
        <div className={`grow`} />
        <Button
          size={`sm`}
          loading={true}
          onClick={async () => await onChange()}
        >
          {t("admin.system.save")}
        </Button>
      </ParagraphFooter>
    </Paragraph>
  );
}

function Search({ data, dispatch, onChange }: CompProps<SearchState>) {
  const { t } = useTranslation();

  const [search, setSearch] = useState<string>("");
  const [searchDialog, setSearchDialog] = useState<boolean>(false);
  const [searchResult, setSearchResult] = useState<string>("");
  const [searchLoading, setSearchLoading] = useState<boolean>(false);

  return (
    <Paragraph
      title={t("admin.system.search")}
      configParagraph={true}
      isCollapsed={true}
    >
      <ParagraphDescription border>
        {t("admin.system.searchTip")}
      </ParagraphDescription>
      <ParagraphItem>
        <Label>{t("admin.system.searchProvider")}</Label>
        <Select
          value={data.provider || "auto"}
          onValueChange={(value) =>
            dispatch({ type: "update:search.provider", value })
          }
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="auto">{t("admin.system.searchAutomatic")}</SelectItem>
            <SelectItem value="parallel">Parallel Search MCP</SelectItem>
            <SelectItem value="searxng">SearXNG</SelectItem>
          </SelectContent>
        </Select>
      </ParagraphItem>
      <ParagraphDescription border>
        {t("admin.system.searchProviderTip")}
      </ParagraphDescription>
      <ParagraphItem>
        <Label>{t("admin.system.searchEndpoint")}</Label>
        <Input
          value={data.endpoint}
          onChange={(e) =>
            dispatch({
              type: "update:search.endpoint",
              value: e.target.value,
            })
          }
          placeholder={t("admin.system.searchPlaceholder")}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.searchEngines")}</Label>
        <MultiCombobox
          value={data.engines}
          onChange={(value) => {
            dispatch({ type: "update:search.engines", value });
          }}
          list={[
            "google",
            "bing",
            "duckduckgo",
            "qwant",
            "brave",
            "mojeek",
            "arxiv",
            "crossref",
            "youtube",
            "bilibili",
            "presearch",
            "yahoo",
            "wiby",
            "seznam",
            "goo",
            "naver",
            "wikidata",
            "wikipedia",
            "wikimini",
            "wikibooks",
            "wikiquote",
            "wikisource",
            "wikispecies",
            "wikiversity",
            "wikivoyage",
            "ask",
            "currency",
            "yep",
            "yacy",
            "genius",
            "github",
            "gitlab",
            "gitea.com",
            "bitbucket",
            "codeberg",
            "mdn",
          ]}
          placeholder={t("admin.system.searchEnginesPlaceholder", {
            length: (data.engines || []).length,
          })}
          searchPlaceholder={t("admin.system.searchEnginesSearchPlaceholder")}
        />
      </ParagraphItem>
      {data.engines.length === 0 && (
        <ParagraphDescription border>
          {t("admin.system.searchEnginesEmptyTip")}
        </ParagraphDescription>
      )}
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.searchImageProxy")}
          <Tips content={t("admin.system.searchImageProxyTip")} />
        </Label>
        <Switch
          checked={data.image_proxy}
          onCheckedChange={(value) => {
            dispatch({ type: "update:search.image_proxy", value });
          }}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label className={`flex flex-row items-center`}>
          {t("admin.system.searchCrop")}
          <Tips content={t("admin.system.searchCropTip")} />
        </Label>
        <Switch
          checked={data.crop}
          onCheckedChange={(value) => {
            dispatch({ type: "update:search.crop", value });
          }}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.searchCropLen")}</Label>
        <NumberInput
          value={data.crop_len}
          onValueChange={(value) =>
            dispatch({ type: "update:search.crop_len", value })
          }
          min={1}
          disabled={!data.crop}
        />
      </ParagraphItem>
      <ParagraphItem>
        <Label>{t("admin.system.searchSafeSearch")}</Label>
        <Combobox
          value={["none", "moderation", "strict"][data.safe_search] || "none"}
          onChange={(value) => {
            dispatch({
              type: "update:search.safe_search",
              value: ["none", "moderation", "strict"].indexOf(value),
            });
          }}
          list={["none", "moderation", "strict"]}
          listTranslated={`admin.system.searchSafeSearchModes`}
          hideSearchBar
        />
      </ParagraphItem>
      <ParagraphFooter>
        <div className={`grow`} />
        <Dialog open={searchDialog} onOpenChange={setSearchDialog}>
          <DialogTrigger asChild>
            <Button variant={`outline`} size={`sm`}>
              {t("admin.system.searchTest")}
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("admin.system.searchTest")}</DialogTitle>
              <FlexibleTextarea
                placeholder={t("admin.system.searchTestTip")}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
              {(searchLoading || searchResult) && (
                <div
                  className={`mt-2 border rounded-md p-4 flex items-center justify-center flex-col`}
                >
                  {searchLoading ? (
                    <Loader2 className={`h-4 w-4 animate-spin`} />
                  ) : (
                    <>
                      <p className={`text-sm mb-1`}>{t("admin.system.searchResult")}</p>
                      <Textarea value={searchResult} rows={5} readOnly />
                    </>
                  )}
                </div>
              )}
            </DialogHeader>
            <DialogFooter>
              <Button
                variant={`outline`}
                onClick={() => {
                  setSearch("");
                  setSearchDialog(false);
                }}
              >
                {t("admin.cancel")}
              </Button>
              <Button
                variant={`default`}
                loading={true}
                onClick={async () => {
                  await onChange();

                  setSearchResult("");
                  setSearchLoading(true);
                  const res = await testWebSearching(search);
                  if (res.status) setSearchResult(res.result);

                  withNotify(t, res, true);
                  setSearchLoading(false);
                }}
              >
                {t("admin.confirm")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        <Button
          size={`sm`}
          loading={true}
          onClick={async () => await onChange()}
        >
          {t("admin.system.save")}
        </Button>
      </ParagraphFooter>
    </Paragraph>
  );
}

function System() {
  const { t } = useTranslation();
  const [data, setData] = useReducer(
    formReducer<SystemProps>(),
    initialSystemState,
  );

  const [loading, setLoading] = useState<boolean>(false);

  const doSaving = async (doToast?: boolean) => {
    const res = await setConfig(data);

    if (doToast !== false) withNotify(t, res, true);
  };

  const doRefresh = async () => {
    setLoading(true);
    const res = await getConfig();
    setLoading(false);
    withNotify(t, res);
    if (res.status) {
      setData({ type: "set", value: res.data });
    }
  };

  useEffectAsync(doRefresh, []);

  return (
    <div className={`system`}>
      <Card className={`admin-card system-card`}>
        <CardHeader className={`select-none`}>
          <CardTitle>{t("admin.settings")}</CardTitle>
        </CardHeader>
        <CardContent className={`flex flex-col gap-1`}>
          <div className={`system-actions flex flex-row`}>
            <div className={`grow`} />
            <Button
              size={`icon`}
              variant={`outline`}
              loading={true}
              className={`mr-2`}
              onClick={async () => await doRefresh()}
            >
              <RotateCw className={cn(loading && `animate-spin`, `h-4 w-4`)} />
            </Button>
            <Button
              size={`icon`}
              loading={true}
              onClick={async () => await doSaving()}
            >
              <Save className={`h-4 w-4`} />
            </Button>
          </div>
          <General
            form={data}
            data={data.general}
            dispatch={setData}
            onChange={doSaving}
          />
          <Site
            form={data}
            data={data.site}
            dispatch={setData}
            onChange={doSaving}
          />
          <Mail
            form={data}
            data={data.mail}
            dispatch={setData}
            onChange={doSaving}
          />
          <Search
            form={data}
            data={data.search}
            dispatch={setData}
            onChange={doSaving}
          />
          <Common
            form={data}
            data={data.common}
            dispatch={setData}
            onChange={doSaving}
          />
        </CardContent>
      </Card>
    </div>
  );
}

export default System;
