import React, {Fragment, useEffect, useState} from "react";
import "primeicons/primeicons.css";
import Sidebar from "./Sidebar";
import Chat from "./Chat";
import {EventsOn} from "../wailsjs/runtime";
import { Transition } from "@headlessui/react";
import {registerTheme} from "echarts";

const App = () => {
    registerTheme('my_theme', {
        "textStyle": {
            "color": "#cccccc",
        },
        "title": {
            "textStyle": {
                "color": "#cccccc",
            }
        }
    });

    const [curConversationID, setCurConversationID] = useState<number | null>(
        null
    );
    const [errorMessage, setErrorMessage] = useState("");

    useEffect(() => {
        return EventsOn("async-error", (err: string) => {
            setErrorMessage(err);
        });
    }, []);

    useEffect(() => {
        let timeout: number;
        if (errorMessage) {
            timeout = setTimeout(() => {
                setErrorMessage("");
            }, 5000);
        }
        return () => clearTimeout(timeout);
    }, [errorMessage]);

    return (
        <div className="flex h-screen overflow-hidden">
            <Sidebar curConversationID={curConversationID} setCurConversationID={setCurConversationID}/>
            <Chat conversationID={curConversationID} setConversationID={setCurConversationID}/>
            <Transition
                show={errorMessage != ""}
                as={Fragment}
                enter="ease-out duration-300"
                enterFrom="opacity-0 scale-95"
                enterTo="opacity-100 scale-100"
                leave="ease-in duration-200"
                leaveFrom="opacity-100 scale-100"
                leaveTo="opacity-0 scale-95"
            >
                <div className="fixed right-4 bottom-4 z-50 rounded-md bg-red-400 text-white p-2">
                    {errorMessage}
                </div>
            </Transition>
        </div>
    );
}

export default App;
