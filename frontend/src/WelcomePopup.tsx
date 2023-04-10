import {Settings} from "iconoir-react";
import React, {Fragment, useEffect, useState} from "react";
import {Dialog, Listbox, Switch, Transition} from "@headlessui/react";
import {GetSettings, SaveSettings} from "../wailsjs/go/main/App";
import {database} from "../wailsjs/go/models";
import {BrowserOpenURL} from "../wailsjs/runtime";

interface Props {
    className?: string;
}

const WelcomePopup = ({className}: Props) => {
    const [settings, setSettings] = useState<database.Settings>();
    const [openAiApiKey, setOpenAiApiKey] = useState("");

    useEffect(() => {
        GetSettings().then((curSettings) => {
            setSettings(curSettings);
            setOpenAiApiKey(curSettings.openAiApiKey);
        });
    }, []);

    const saveSettings = async () => {
        let settings = await GetSettings();
        settings.openAiApiKey = openAiApiKey;
        let newSettings = await SaveSettings(settings);
        setSettings(newSettings);
    }

    if (!settings) {
        return <></>
    }

    return (
        <>
            <Transition show={!settings?.openAiApiKey} as={Fragment}>
                <Dialog open={!settings?.openAiApiKey} onClose={() => {}}>
                    <Transition.Child
                        as={Fragment}
                        enter="ease-out duration-300"
                        enterFrom="opacity-0"
                        enterTo="opacity-50"
                        leave="ease-in duration-200"
                        leaveFrom="opacity-50"
                        leaveTo="opacity-0"
                    >
                        <div
                            className="fixed inset-0 bg-gray-800 opacity-50 z-30"
                        ></div>
                    </Transition.Child>
                    <Transition.Child
                        as={Fragment}
                        enter="ease-out duration-300"
                        enterFrom="opacity-0 scale-95"
                        enterTo="opacity-100 scale-100"
                        leave="ease-in duration-200"
                        leaveFrom="opacity-100 scale-100"
                        leaveTo="opacity-0 scale-95"
                    >
                        <Dialog.Panel
                            className="flex flex-col fixed inset-60 z-40 bg-gray-900 rounded-md p-4 overflow-hidden">
                            <Dialog.Title className="text-lg font-bold text-gray-400 mb-4">Welcome!</Dialog.Title>
                            <p className="text-gray-400 p-2">To start with, please provide your OpenAI API key. It's required to do anything in this app. You can generate your API key on the <b className="text-blue-300 cursor-pointer" onClick={() => BrowserOpenURL("https://platform.openai.com")}>OpenAI developer platform</b>.</p>
                            <div className="divide-y divide-gray-700 h-full overflow-y-auto">
                                <div className="flex items-center justify-between p-2">
                                    <p className="text-gray-400">OpenAI API Key</p>
                                    <input type="password"
                                           value={openAiApiKey}
                                           onChange={(event) => setOpenAiApiKey(event.target.value)}
                                           className="border border-gray-300 border-opacity-50 p-2 h-8 bg-gray-700 text-gray-300 rounded-md"/>
                                </div>
                            </div>
                            <div className="flex justify-end">
                                <button
                                    type="button"
                                    onClick={async () => await saveSettings()}
                                    className={`${openAiApiKey !== "" ? "bg-blue-500" : "bg-gray-500"} text-white p-2 rounded-md mt-2`}
                                >
                                    Submit
                                </button>
                            </div>
                        </Dialog.Panel>
                    </Transition.Child>
                </Dialog>
            </Transition>
        </>
    )
}

export default WelcomePopup;
