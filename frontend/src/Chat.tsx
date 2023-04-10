import {CancelGeneration, GetConversation, Messages, RerunFromMessage, SendMessage} from "../wailsjs/go/main/App";
import React, {useEffect, useRef, useState} from "react";
import {database} from "../wailsjs/go/models";
import {EventsOn} from "../wailsjs/runtime";
import ReactMarkdown from "react-markdown";
import {Prism as SyntaxHighlighter} from "react-syntax-highlighter";
import {dracula} from "react-syntax-highlighter/dist/esm/styles/prism";
import {EditPencil, MinusCircle, RefreshDouble} from "iconoir-react";
import Message = database.Message;
import Conversation = database.Conversation;
import {capitalizeFirstLetter, isJSONString} from "./helpers";
import ChatInputForm from "./ChatInputForm";
import ReactECharts from 'echarts-for-react';
import MessageBubble from "./Message";

interface Props {
    conversationID: number | null;
    setConversationID: (conversationID: number) => void;
}

const Chat = ({conversationID, setConversationID}: Props) => {
    const [messages, setMessages] = useState<Array<Message>>([]);
    const [curConversation, setCurConversation] = useState<database.Conversation | null>(null);
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (conversationID === null) {
            setMessages([]);
            setCurConversation(null);
            return;
        }
        Messages(conversationID).then((messages) => {
            setMessages(messages);
        });
        GetConversation(conversationID).then((conversation: Conversation) => {
            setCurConversation(conversation);
        })
    }, [conversationID]);

    useEffect(() => {
        if (conversationID === null) {
            setMessages([]);
            setCurConversation(null);
            return;
        }
        return EventsOn(`conversation-${conversationID}-updated`, (data: any) => {
            Messages(conversationID).then((messages) => {
                setMessages(messages);
            });
            GetConversation(conversationID).then((conversation: Conversation) => {
                setCurConversation(conversation);
            })
        })
    }, [conversationID]);

    useEffect(() => {
        if (messagesContainerRef.current) {
            messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
        }
    }, [messages]); // TODO: Change dependency from messages to something like "sent message" so it doesn't happen on assistant responses.

    return <div className="w-3/4 bg-gray-800">
        <div className="flex flex-col h-full">
            <div className="relative flex-1 overflow-hidden">
                <div
                    ref={messagesContainerRef}
                    className="flex-1 h-full overflow-y-auto px-4 py-2"
                >
                    {messages.map((message, index) => (
                        <div key={index}
                             className={`flex flex-col ${message.author == 'user' ? "items-end" : "items-start"}`}>
                            <MessageBubble message={message}/>
                        </div>
                    ))}
                </div>
                {curConversation?.generating &&
                  <MinusCircle className="absolute left-4 bottom-1 hover:text-gray-400" onClick={async () => {
                      if (curConversation !== null) {
                          await CancelGeneration(curConversation.id)
                      }
                  }}/>}
            </div>
            <ChatInputForm disabled={curConversation?.generating || false} conversationID={conversationID} setConversationID={setConversationID}/>
        </div>
    </div>
}

export default Chat;
