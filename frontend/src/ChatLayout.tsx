import React, {Fragment, useEffect, useRef, useState} from "react";
import ReactMarkdown from "react-markdown";
import {Light as SyntaxHighlighter} from "react-syntax-highlighter";
import {dracula} from "react-syntax-highlighter/dist/esm/styles/hljs";
import "highlight.js/styles/github-dark-dimmed.css";
import {Settings} from 'iconoir-react';
import {Switch, Transition} from '@headlessui/react'

type ChatMessage = {
    message: string;
    isSentByMe: boolean;
};

const ChatLayout = () => {
    const [messages, setMessages] = useState<Array<ChatMessage>>([]);
    const [inputText, setInputText] = useState("");
    const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);
    const [openAIApiKey, setOpenAIApiKey] = useState("");
    const [toggleOption1, setToggleOption1] = useState(false);
    const [textInputOption1, setTextInputOption1] = useState("");
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (messagesContainerRef.current) {
            messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
        }
    }, [messages]);

    const handleKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            handleSubmit();
        }
    };

    const handleSubmit = (debugMode: boolean = false) => {
        if (inputText.trim() !== "") {
            setMessages([
                ...messages,
                {message: inputText, isSentByMe: !debugMode},
            ]);
            setInputText("");
        }
    };

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

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        alert("Code copied to clipboard!");
    };

    const renderMarkdown = (message: ChatMessage) => {
        return (
            <ReactMarkdown
                children={message.message}
                components={{
                    code({node, inline, className, children, ...props}) {
                        const match = /language-(\w+)/.exec(className || "");
                        if (inline) {
                            return (
                                <code className={className} {...props}>
                                    {children}
                                </code>
                            );
                        }
                        return (
                            <div
                                className="relative"
                                onMouseEnter={(e) => {
                                    const button = e.currentTarget.querySelector("button");
                                    if (button !== null) {
                                        button.style.opacity = "100";
                                    }
                                }}
                                onMouseLeave={(e) => {
                                    const button = e.currentTarget.querySelector("button");
                                    if (button !== null) {
                                        button.style.opacity = "0";
                                    }
                                }}
                            >
                                <button
                                    className="absolute top-0 right-0 bg-gray-700 text-white py-1 px-2 rounded-md opacity-0 hover:opacity-100 transition-opacity"
                                    onClick={() => copyToClipboard(children as string)}
                                >
                                    Copy
                                </button>
                                {match ? (
                                    <SyntaxHighlighter
                                        className="rounded-md"
                                        children={String(children).replace(/\n$/, "")}
                                        style={dracula as any}
                                        language={match[1]}
                                        PreTag="div"
                                        {...props}
                                    />
                                ) : (
                                    <SyntaxHighlighter
                                        className="rounded-md"
                                        children={String(children).replace(/\n$/, "")}
                                        style={dracula as any}
                                        PreTag="div"
                                        {...props}
                                    />
                                )}
                            </div>
                        );
                    },
                }}
                className={`bg-${message.isSentByMe ? "green" : "gray"}-700 py-2 px-4 rounded-md text-white inline-block max-w-full`}
            />
        );
    };

    return (
        <div className="flex h-screen overflow-hidden">
            <div className="w-1/4 border-r border-gray-300 border-opacity-50 p-4 bg-gray-900">
                <h2 className="font-bold text-lg mb-4 text-gray-300">Conversations</h2>
                {conversations.map((conversation, index) => (
                    <div
                        key={conversation.id}
                        className="flex items-center mb-4 cursor-pointer border-b border-gray-400 py-2"
                    >
                        <div className="w-10 h-10 rounded-full bg-gray-300 mr-2"></div>
                        <div className="flex-1 text-gray-500">
                            <div className="flex justify-between">
                                <p className="text-sm">{conversation.timestamp}</p>
                            </div>
                            <p className="text-gray-500">{conversation.message}</p>
                        </div>
                    </div>
                ))}
                <div onClick={() => setIsSettingsModalOpen(true)} className="absolute bottom-4 left-4 cursor-pointer">
                    <Settings className="text-gray-500"/>
                </div>
            </div>
            <div className="w-3/4 bg-gray-800">
                <div className="flex flex-col h-full">
                    <div
                        ref={messagesContainerRef}
                        className="flex-1 overflow-y-auto px-4 py-2"
                    >
                        {messages.map((message, index) => (
                            <div key={index}
                                 className={`mb-4 flex flex-col ${message.isSentByMe ? "items-end" : "items-start"}`}>
                                <div className="text-gray-500 mb-1">
                                    {new Date().toLocaleString()}
                                </div>
                                <div
                                    className={`py-2 px-4 rounded-md text-white inline-block relative max-w-full`}
                                    style={{wordWrap: "break-word"}}
                                >
                                    {renderMarkdown(message)}
                                </div>
                            </div>
                        ))}
                    </div>
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
                        <div className="flex justify-between">
                            <button
                                type="button"
                                onClick={() => handleSubmit(false)}
                                className="bg-blue-500 text-white p-2 rounded-md mt-2"
                            >
                                Send
                            </button>
                            <button
                                type="button"
                                onClick={() => handleSubmit(true)}
                                className="bg-red-500 text-white p-2 rounded-md mt-2"
                            >
                                Send (Other, Debug)
                            </button>
                        </div>
                    </form>
                </div>
            </div>
            {isSettingsModalOpen && (
                // Change to HeadlessUI Dialog
                <>
                    {/* Background overlay */}
                    <div
                        onClick={() => setIsSettingsModalOpen(false)}
                        className="fixed inset-0 bg-gray-800 opacity-50 z-30"
                    ></div>

                    {/* Settings modal */}
                    <Transition
                        show={isSettingsModalOpen}
                        enter="transition duration-100 ease-out"
                        enterFrom="transform scale-95 opacity-0"
                        enterTo="transform scale-100 opacity-100"
                        leave="transition duration-75 ease-out"
                        leaveFrom="transform scale-100 opacity-100"
                        leaveTo="transform scale-95 opacity-0"
                        as={Fragment}
                    >
                        <div className="fixed inset-40 z-40 bg-gray-900 rounded-md p-4">
                            <h2 className="text-lg font-bold text-gray-400 mb-4">Settings</h2>

                            {/* Settings list */}
                            <div className="divide-y divide-gray-700 max-h-80 overflow-y-auto">
                                <div className="flex items-center justify-between p-2">
                                    <p className="text-gray-400">OpenAI API Key</p>
                                    {/*<Password inputStyle="text-gray-300" value={openAIApiKey} onChange={(e) => setOpenAIApiKey(e.target.value)} feedback={false}/>*/}

                                    <input type="password"
                                           className="border border-gray-300 border-opacity-50 p-2 h-8 bg-gray-700 text-gray-300 rounded-md"/>
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
                                        onKeyDown={handleKeyDown}
                                        className="border border-gray-300 border-opacity-50 p-2 w-full h-32 bg-gray-700 text-gray-300 resize-none rounded-md"
                                    />
                                </div>

                                {/* Add more options here */}

                            </div>
                        </div>
                    </Transition>
                </>
            )}
        </div>
    );
};

export default ChatLayout;

// TODO: Maybe just post a skeleton to chatgpt and ask it to add the settings modal?
