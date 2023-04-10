import React, {useEffect, useState} from "react";
import {database} from "../wailsjs/go/models";
import AppSettingsButton from "./AppSettingsButton";
import {Conversations, DeleteConversation} from "../wailsjs/go/main/App";
import {EventsOn} from "../wailsjs/runtime";
import {Bin, EditPencil} from "iconoir-react";
import Conversation = database.Conversation;
import ConversationSettingsButton from "./ConversationSettingsButton";

interface Props {
    curConversationID: number | null;
    setCurConversationID: (conversationID: number | null) => void;
}

const Sidebar = ({curConversationID, setCurConversationID}: Props) => {
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

    const onConversationDelete = async (id: number) => {
        await DeleteConversation(id);
        setCurConversationID(null);
    }

    return (
        <div className="flex flex-col h-full w-1/4 border-r border-gray-300 border-opacity-50 bg-gray-900 p-2">
            <div className="flex-1 overflow-hidden flex flex-col w-full border rounded-md border-gray-300 border-opacity-50">
                <h2 className="font-bold text-lg text-gray-300 p-3">Conversations</h2>
                <div className="overflow-y-auto divide-y divide-gray-700 border-t border-gray-300 border-opacity-50">
                    <div
                        className={`flex items-center cursor-pointer py-2 hover:bg-gray-700`}
                        onClick={() => {
                            setCurConversationID(null)
                        }}
                    >
                        {/*<div className="w-10 h-10 rounded-full bg-gray-300 mr-2"></div>*/}
                        <div className="relative flex-1 text-gray-500 px-2">
                            <p className="text-gray-500">Create new conversation...</p>
                            <ConversationSettingsButton className="absolute scale-75 top-0 right-1 text-gray-500 hover:text-gray-400" conversationSettingsID={null}/>
                        </div>
                    </div>
                    {conversations.map((conversation, index) => (
                        <div
                            key={conversation.id}
                            className={`relative flex items-center cursor-pointer py-2 ${curConversationID == conversation.id ? "bg-gray-800" : ""} hover:bg-gray-700`}
                            onClick={() => {
                                setCurConversationID(conversation.id)
                            }}
                        >
                            {/*<div className="w-10 h-10 rounded-full bg-gray-300 mr-2"></div>*/}
                            <div className="flex-1 text-gray-500 px-2">
                                <div className="flex justify-between">
                                    <p className="text-sm">Today</p>
                                </div>
                                <p className="text-gray-500">{conversation.title}</p>
                            </div>
                            <ConversationSettingsButton className="absolute scale-75 top-1 right-7 text-gray-500 hover:text-gray-400" conversationSettingsID={conversation.conversationSettingsID}/>
                            <Bin className="absolute scale-75 top-1 right-1 text-gray-500 hover:text-red-400"
                                 onClick={async () => onConversationDelete(conversation.id)}/>
                        </div>
                    ))}
                </div>
            </div>
            <div className="h-12"></div>
            <AppSettingsButton className="absolute bottom-4 left-4"/>
        </div>
    )
}

export default Sidebar
