import {Settings} from "iconoir-react";
import React from "react";
import {database} from "../wailsjs/go/models";
import Conversation = database.Conversation;
import SettingsButton from "./SettingsButton";

interface Props {
}

const Sidebar = ({}: Props) => {
    const conversations = [
        {
            id: 1,
            message: "Hey, how's it going?",
            timestamp: "10:23 AM",
        },
        {
            id: 2,
            message: "Can you send me the report?",
            timestamp: "Yesterday",
        },
        {
            id: 3,
            message: "I'll be late to the meeting",
            timestamp: "Tuesday",
        },
    ];

    return (
        <div className="w-1/4 border-r border-gray-300 border-opacity-50 p-4 bg-gray-900">
            <h2 className="font-bold text-lg mb-4 text-gray-300">Conversations</h2>
            {conversations.map((conversation, index) => (
                <div
                    key={conversation.id}
                    className="flex items-center mb-4 cursor-pointer border-b border-gray-400 py-2"
                >
                    {/*<div className="w-10 h-10 rounded-full bg-gray-300 mr-2"></div>*/}
                    <div className="flex-1 text-gray-500">
                        <div className="flex justify-between">
                            <p className="text-sm">Today</p>
                        </div>
                        <p className="text-gray-500">{conversation.message}</p>
                    </div>
                </div>
            ))}
            <SettingsButton/>
        </div>
    )
}

export default Sidebar
