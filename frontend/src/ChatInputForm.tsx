import {Settings} from "iconoir-react";
import React, {Fragment, useEffect, useState} from "react";
import {Dialog, Listbox, Transition} from "@headlessui/react";
import {GetSettings, SaveSettings, SendMessage} from "../wailsjs/go/main/App";
import {database} from "../wailsjs/go/models";

interface Props {
    disabled: boolean;
    conversationID: number;
    setConversationID: (conversationID: number) => void;
}

const ChatInputForm = ({disabled, conversationID, setConversationID}: Props) => {
    const [inputText, setInputText] = useState("");

    const handleKeyDown = async (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            await handleSubmit();
        }
    };

    const handleSubmit = async () => {
        if (disabled) {
            return;
        }
        if (inputText.trim() !== "") {
            let message = await SendMessage(conversationID !== null ? conversationID : -1, inputText);
            setInputText("");
            setConversationID(message.conversationID);
        }
    };

    return (
        <form
            onSubmit={(e) => e.preventDefault()}
            className="flex flex-col h-48 px-4 py-2"
        >
                <textarea
                    value={inputText}
                    onChange={(event) => setInputText(event.target.value)}
                    onKeyDown={handleKeyDown}
                    className="border border-gray-300 border-opacity-50 p-2 w-full h-32 bg-gray-900 text-white resize-none rounded-md"
                />
            <div className="flex justify-end">
                <button
                    type="button"
                    onClick={handleSubmit}
                    className={`${disabled ? "bg-gray-500" : "bg-blue-500"} text-white p-2 rounded-md mt-2`}
                >
                    Send
                </button>
            </div>
        </form>
    )
}

export default ChatInputForm;
