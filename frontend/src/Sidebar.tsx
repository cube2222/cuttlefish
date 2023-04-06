import React, {useEffect, useState} from "react";
import {database} from "../wailsjs/go/models";
import SettingsButton from "./SettingsButton";
import {Conversations} from "../wailsjs/go/main/App";
import {EventsOn} from "../wailsjs/runtime";
import Conversation = database.Conversation;

interface Props {
}

const Sidebar = ({}: Props) => {
    const [conversations, setConversations] = useState<Array<Conversation>>([]);

    useEffect(() => {
        Conversations().then((conversations) => {
            setConversations(conversations);
        });
    }, []);

    useEffect(() => {
        return EventsOn(`conversations-updated`, (data: any) => {
            Conversations().then((conversations) => {
                setConversations(conversations);
            });
        })
    }, []);


    return (
        <div className="w-1/4 border-r border-gray-300 border-opacity-50 p-4 bg-gray-900">
            <h2 className="font-bold text-lg mb-4 text-gray-300">Conversations</h2>
            <div className="h-5/6 overflow-y-auto">
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
                            <p className="text-gray-500">{conversation.title}</p>
                        </div>
                    </div>
                ))}
            </div>
            <SettingsButton/>
        </div>
    )
}

export default Sidebar
