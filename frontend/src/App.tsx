import {useState} from "react";
import {Messages} from "../wailsjs/go/main/App";
import "./ChatLayout";
import ChatLayout from "./ChatLayout";
import "primeicons/primeicons.css";

function App() {
    return <ChatLayout/>;
    // return (
    //     <div className="min-h-screen bg-white grid grid-cols-1 place-items-center justify-items-center mx-auto py-8">
    //         <div className="text-blue-900 text-2xl font-bold font-mono">
    //             <h1 className="content-center">{text}</h1>
    //         </div>
    //         <div className="w-fit max-w-md">
    //             <button className="btn" onClick={(t) => Greet('Kuba').then((res) => updateText(res))}>Greet</button>
    //         </div>
    //     </div>
    // )
}

export default App
