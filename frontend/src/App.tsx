import React from "react";
import "primeicons/primeicons.css";
import Sidebar from "./Sidebar";
import Chat from "./Chat";

function App() {
    return (
        <div className="flex h-screen overflow-hidden">
            <Sidebar/>
            <Chat/>
        </div>
    )
}

export default App
