import {Approve, CancelGeneration, GetConversation, ListApprovalRequests, Messages} from "../wailsjs/go/main/App";
import React, {useEffect, useRef, useState} from "react";
import {database, main} from "../wailsjs/go/models";
import {EventsOn} from "../wailsjs/runtime";
import {MinusCircle} from "iconoir-react";
import ChatInputForm from "./ChatInputForm";
import MessageBubble from "./Message";
import Message = database.Message;
import Conversation = database.Conversation;
import ApprovalRequest = main.ApprovalRequest;

interface Props {
    conversationID: number | null;
    setConversationID: (conversationID: number) => void;
}

const Chat = ({conversationID, setConversationID}: Props) => {
    const [messages, setMessages] = useState<Array<Message>>([]);
    const [approvalRequests, setApprovalRequests] = useState<Array<ApprovalRequest>>([]);
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
        ListApprovalRequests(conversationID).then((requests) => {
            setApprovalRequests(requests);
        })
    }, [conversationID]);

    useEffect(() => {
        if (conversationID === null) {
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
        if (conversationID === null) {
            return;
        }
        return EventsOn(`conversation-${conversationID}-approvals-updated`, (data: any) => {
            ListApprovalRequests(conversationID).then((requests) => {
                setApprovalRequests(requests);
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
                  <MinusCircle className="absolute left-4 bottom-1 hover:text-gray-400 cursor-pointer" onClick={async () => {
                      if (curConversation !== null) {
                          await CancelGeneration(curConversation.id)
                      }
                  }}/>}
                {conversationID && approvalRequests.map((request, index) => {
                    // There's at most one.
                    return <div className="absolute left-12 bottom-1 bg-green-400 hover:bg-green-300 opacity-75 text-gray-800 rounded-full px-2 cursor-pointer" onClick={async () => await Approve(conversationID, request.id)}>
                        Click to approve: "{request.message}"
                    </div>
                })}
            </div>
            <ChatInputForm disabled={curConversation?.generating || false} conversationID={conversationID}
                           setConversationID={setConversationID}/>
        </div>
    </div>
}

export default Chat;
