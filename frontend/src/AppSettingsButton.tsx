import {Settings} from "iconoir-react";
import React, {Fragment, useEffect, useState} from "react";
import {Dialog, Listbox, Switch, Transition} from "@headlessui/react";
import {GetSettings, SaveSettings} from "../wailsjs/go/main/App";
import {main} from "../wailsjs/go/models";

interface Props {
    className?: string;
}

const AppSettingsButton = ({className}: Props) => {
    const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);
    const [settings, setSettings] = useState<main.Settings>();
    const [openAiApiKey, setOpenAiApiKey] = useState("");
    const [toggleOption1, setToggleOption1] = useState(false);
    const [textInputOption1, setTextInputOption1] = useState("");
    const [model, setModel] = useState("gpt-3.5-turbo");
    const [changed, setChanged] = useState(false);

    useEffect(() => {
        GetSettings().then((curSettings) => {
            setSettings(curSettings);
            setOpenAiApiKey(curSettings.openAiApiKey);
            setModel(curSettings.model);
        });
    }, []);

    useEffect(() => {
        if (!settings) {
            return;
        }
        setChanged(
            openAiApiKey !== settings.openAiApiKey
            || model !== settings.model
        );
    }, [openAiApiKey, model])

    const saveSettings = async () => {
        await SaveSettings({
            openAiApiKey: openAiApiKey,
            model: model,
        });
        setChanged(false);
    }

    return (
        <>
            <div onClick={() => setIsSettingsModalOpen(true)} className={"cursor-pointer " + className}>
                <Settings className="text-gray-500 hover:text-gray-400"/>
            </div>
            <Transition show={isSettingsModalOpen} as={Fragment}>
                <Dialog open={isSettingsModalOpen} onClose={() => setIsSettingsModalOpen(false)}>
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
                            onClick={() => setIsSettingsModalOpen(false)}
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
                        <Dialog.Panel className="fixed inset-40 z-40 bg-gray-900 rounded-md p-4">
                            <Dialog.Title className="text-lg font-bold text-gray-400 mb-4">Settings</Dialog.Title>
                            <div className="divide-y divide-gray-700 max-h-80 overflow-y-auto">
                                <div className="flex items-center justify-between p-2">
                                    <p className="text-gray-400">OpenAI API Key</p>
                                    <input type="password"
                                           value={openAiApiKey}
                                           onChange={(event) => setOpenAiApiKey(event.target.value)}
                                           className="border border-gray-300 border-opacity-50 p-2 h-8 bg-gray-700 text-gray-300 rounded-md"/>
                                </div>
                                <div className="flex items-center justify-between p-2">
                                    <p className="text-gray-400">Model</p>
                                    <div className="w-1/3 max-w-xs">
                                        <Listbox value={model} onChange={setModel}>
                                            <div className="relative">
                                                <Listbox.Button
                                                    className="duration-150 cursor-default relative w-full border border-gray-300 border-opacity-50 rounded-md bg-gray-700 text-gray-300 pl-3 py-1.5 text-left hover:bg-gray-600">{model}</Listbox.Button>
                                                <Listbox.Options
                                                    className="bg-gray-700 absolute mt-1 w-full rounded-md bg-white shadow-lg max-h-60 rounded-md z-40 divide-y divide-gray-600">
                                                    {["gpt-3.5-turbo", "gpt-4"].map((model) => (
                                                        <Listbox.Option
                                                            key={model}
                                                            value={model}
                                                            className="duration-150 text-gray-300 cursor-default pl-4 py-2 rounded-md hover:bg-gray-600"
                                                        >
                                                    <span className="block truncate">
                                                        {model}
                                                    </span>
                                                        </Listbox.Option>
                                                    ))}
                                                </Listbox.Options>
                                            </div>
                                        </Listbox>
                                    </div>
                                </div>
                                <div className="flex items-center justify-between p-2">
                                    <p className="text-gray-400">Toggle Option 1</p>
                                    <Switch
                                        checked={toggleOption1}
                                        onChange={setToggleOption1}
                                        className={`${
                                            toggleOption1 ? 'bg-gray-400' : 'bg-gray-700'
                                        } relative inline-flex h-6 w-11 items-center rounded-full border border-gray-300 border-opacity-50`}
                                    >
                                    <span
                                        className={`${
                                            toggleOption1 ? 'translate-x-6' : 'translate-x-1'
                                        } inline-block h-4 w-4 transform rounded-full bg-gray-200 transition`}
                                    />
                                    </Switch>
                                </div>
                                <div className="flex flex-col p-2">
                                    <label htmlFor="textInput1" className="text-gray-400 mb-1">
                                        Text Input Option
                                    </label>
                                    <textarea
                                        value={textInputOption1}
                                        onChange={(event) => setTextInputOption1(event.target.value)}
                                        // onKeyDown={handleKeyDown}
                                        className="border border-gray-300 border-opacity-50 p-2 w-full h-32 bg-gray-700 text-gray-300 resize-none rounded-md"
                                    />
                                </div>
                            </div>
                            <div className="flex justify-end">
                                <button
                                    type="button"
                                    onClick={async () => await saveSettings()}
                                    className={`${changed ? "bg-blue-500" : "bg-gray-500"} text-white p-2 rounded-md mt-2`}
                                >
                                    Save
                                </button>
                            </div>
                        </Dialog.Panel>
                    </Transition.Child>
                </Dialog>
            </Transition>
        </>
    )
}

export default AppSettingsButton;
