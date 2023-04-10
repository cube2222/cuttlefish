import {EditPencil} from "iconoir-react";
import React, {Fragment, useEffect, useState} from "react";
import {Dialog, Switch, Transition} from "@headlessui/react";
import {
    GetAvailableTools,
    GetConversationSettings,
    GetDefaultConversationSettings, SetDefaultConversationSettings,
    UpdateConversationSettings
} from "../wailsjs/go/main/App";
import {database, main} from "../wailsjs/go/models";
import {capitalizeFirstLetter} from "./helpers";
import AvailableTool = main.AvailableTool;

interface Props {
    className?: string;
    conversationSettingsID: number | null;
}

const ConversationSettingsButton = ({className, conversationSettingsID}: Props) => {
    const isDefault = conversationSettingsID === null;

    const [availableTools, setAvailableTools] = useState<AvailableTool[]>([]);
    const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);
    const [settings, setSettings] = useState<database.ConversationSetting>();
    const [systemPromptTemplate, setSystemPromptTemplate] = useState("");
    const [toolsEnabled, setToolsEnabled] = useState<Set<string>>(new Set());
    const [changed, setChanged] = useState(false);

    useEffect(() => {
        GetAvailableTools().then((tools: AvailableTool[]) => {
            setAvailableTools(tools);
        })
    })

    useEffect(() => {
        if (!isDefault) {
            GetConversationSettings(conversationSettingsID).then((curSettings) => {
                setSettings(curSettings);
                setSystemPromptTemplate(curSettings.systemPromptTemplate);
                setToolsEnabled(new Set(curSettings.toolsEnabled));
            });
        } else {
            GetDefaultConversationSettings().then((curSettings) => {
                setSettings(curSettings);
                setSystemPromptTemplate(curSettings.systemPromptTemplate);
                setToolsEnabled(new Set(curSettings.toolsEnabled));
            });
        }
    }, [isSettingsModalOpen]);

    useEffect(() => {
        if (!settings) {
            return;
        }
        setChanged(
            systemPromptTemplate !== settings.systemPromptTemplate
            || !arraySetsEqual(Array.from(toolsEnabled), settings.toolsEnabled)
        );
    }, [settings, systemPromptTemplate, toolsEnabled])

    const setToolEnabled = (tool: string, enabled: boolean) => {
        let toolsEnabledUpdated = new Set(toolsEnabled);
        if (enabled) {
            toolsEnabledUpdated.add(tool);
        } else {
            toolsEnabledUpdated.delete(tool);
        }
        setToolsEnabled(toolsEnabledUpdated);
    }

    const saveSettings = async () => {
        if (!isDefault) {
            const curSettings = await UpdateConversationSettings({
                id: conversationSettingsID,
                systemPromptTemplate: systemPromptTemplate,
                toolsEnabled: Array.from(toolsEnabled),
            });
            setSettings(curSettings);
        } else {
            const curSettings = await SetDefaultConversationSettings({
                systemPromptTemplate: systemPromptTemplate,
                toolsEnabled: Array.from(toolsEnabled),
            });
            setSettings(curSettings);
        }
        setChanged(false);
    }

    return (
        <>
            <div onClick={() => setIsSettingsModalOpen(true)} className={"cursor-pointer " + className}>
                <EditPencil className="text-gray-500 hover:text-gray-400"/>
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
                        <Dialog.Panel className="flex flex-col fixed inset-40 z-40 bg-gray-900 rounded-md p-4">
                            <Dialog.Title className="text-lg font-bold text-gray-400 mb-4">{isDefault && "Default "}Conversation Settings</Dialog.Title>
                            <div className="flex-1 divide-y divide-gray-700 overflow-y-auto">
                                <div className="flex flex-col p-2">
                                    <label htmlFor="textInput1" className="text-gray-400 mb-2">
                                        System Prompt Template
                                    </label>
                                    <textarea
                                        value={systemPromptTemplate}
                                        onChange={(event) => setSystemPromptTemplate(event.target.value)}
                                        className="border border-gray-300 border-opacity-50 p-2 w-full h-32 bg-gray-700 text-gray-300 resize-none rounded-md"
                                    />
                                </div>
                                <div className="p-2">
                                    <h2 className="text-md font-bold text-gray-400 mb-2">Enabled Tools</h2>
                                    <div className="flex flex-col">
                                        {availableTools.map((tool) => {
                                            return <div className="flex items-center justify-between p-2">
                                                <p className="text-gray-400">{tool.name}</p>
                                                <Switch
                                                    checked={toolsEnabled.has(tool.ID)}
                                                    onChange={(newValue) => setToolEnabled(tool.ID, newValue)}
                                                    className={`${
                                                        toolsEnabled.has(tool.ID) ? 'bg-gray-400' : 'bg-gray-700'
                                                    } relative inline-flex h-6 w-11 items-center rounded-full border border-gray-300 border-opacity-50`}
                                                >
                                                <span
                                                    className={`${
                                                        toolsEnabled.has(tool.ID) ? 'translate-x-6' : 'translate-x-1'
                                                    } inline-block h-4 w-4 transform rounded-full bg-gray-200 transition`}
                                                />
                                                </Switch>
                                            </div>
                                        })}
                                    </div>
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

function arraySetsEqual(arr1: string[], arr2: string[]): boolean {
    const arr1Sorted = arr1.sort();
    const arr2Sorted = arr2.sort();
    if (arr1Sorted.length != arr2Sorted.length) return false;
    return arr1Sorted.every((element, index) => element === arr2Sorted[index]);
}

export default ConversationSettingsButton;
