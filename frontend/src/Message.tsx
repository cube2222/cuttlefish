import {RefreshDouble, Settings} from "iconoir-react";
import React, {Fragment, useEffect, useState} from "react";
import {Dialog, Listbox, Transition} from "@headlessui/react";
import {GetSettings, RerunFromMessage, SaveSettings, SendMessage} from "../wailsjs/go/main/App";
import {database} from "../wailsjs/go/models";
import Message = database.Message;
import {capitalizeFirstLetter, isJSONString} from "./helpers";
import ReactMarkdown from "react-markdown";
import ReactECharts from "echarts-for-react";
import {Prism as SyntaxHighlighter} from "react-syntax-highlighter";
import {dracula} from "react-syntax-highlighter/dist/esm/styles/prism";

interface Props {
    message: Message;
}

const MessageBubble = ({message}: Props) => {
    const [effect, setEffect] = useState(false);

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
    };

    const renderMarkdown = (message: Message) => {
        return (
            // TODO: Custom Thought and Action rendering.
            <ReactMarkdown
                children={message.content}
                components={{
                    code({node, inline, className, children, ...props}) {
                        const match = /language-(\w+)/.exec(className || "");
                        let language = match ? match[1] : null;
                        if (!language && isJSONString(String(children))) {
                            language = "json";
                        }
                        if (language == "action") {
                            // TODO: Render tool use in a special way. I.e. Python should print the python code nicely.
                            language = "json";
                        }
                        if (language == "chart" && isJSONString(String(children))) {
                            return <ReactECharts theme={"my_theme"} style={{ width: "400px", height: "400px" }} option={JSON.parse(String(children))}/>
                        }
                        if (inline) {
                            return (
                                <code className={className} {...props}>
                                    {children}
                                </code>
                            );
                        }
                        return (
                            <div
                                className="relative py-2"
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
                                {language ? (
                                    <SyntaxHighlighter
                                        className="rounded-md"
                                        children={String(children).replace(/\n$/, "")}
                                        style={dracula as any}
                                        language={language}
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
                    img({...props}) {
                        return <img {...props} className="rounded-md w-60 my-2" />
                    }
                }}
                className={`${message.author == 'user' ? "bg-gray-600" : "bg-gray-700"} py-2 px-4 rounded-md ${message.author == 'user' ? "text-gray-200" : "text-gray-300"} inline-block max-w-full`}
            />
        );
    };

    return (
        <div
            className={`flex flex-col max-w-5/6 w-5/6 ${message.author == 'user' ? "items-end" : "items-start"} py-1 px-4 rounded-md text-white inline-block relative`}
            style={{wordWrap: "break-word"}}
        >
            <div className={`${message.author == 'user' ? "text-end" : "text-start"} text-gray-500 p-1 px-2`}>
                {capitalizeFirstLetter(message.author == 'user' ? "you" : message.author)}
            </div>
            <div className="flex flex-row">
                {message.author == 'user' && <RefreshDouble className={`${effect && "animate-refresh_rotate_scaled"} flex-none m-2 scale-75 text-gray-500 hover:text-gray-400 cursor-pointer`} onClick={async () => {
                    setEffect(true);
                    await RerunFromMessage(message.conversationID, message.id);
                }} onAnimationEnd={() => setEffect(false)}></RefreshDouble>}
                {renderMarkdown(message)}
            </div>
        </div>
    )
}

export default MessageBubble;
