<script>
    import { onDestroy } from "svelte/internal";
    import {
        userInfo,
        onlineUserStore,
        chatTabs,
        isTypingStore,
        groupIsTypingStore,
        API_URL,
        IMAGE_URL,
        markGroupMessageAsSeen,
    } from "../../stores";
    import {
        connect,
        sendMessage,
        messages,
        sendDataRequest,
    } from "../../websocket";
    import CloseChat from "../icons/closeChat.svelte";
    import MinimizeChat from "../icons/minimizeChat.svelte";
    import ChatModuleEmojiPicker from "../icons/chatModuleEmojiPicker.svelte";
    import { compute_slots, each } from "svelte/internal";
    import IsTyping from "./isTyping.svelte";
    // Relationship Status
    import { allUsers } from "../../stores";
    import { get } from "svelte/store";
    import Message from "./message.svelte";
    import ChatFollowing from "./chatFollowing.svelte";
    import { selectUser, removeFromActiveChat } from "../../utils";
    export let AvatarPath = "";
    let isOnline = true;
    if (AvatarPath === "") {
        if (isGroup) {
            AvatarPath = "/images/avatars/defaultGroup.png";
        } else AvatarPath = "/images/avatars/default.png";
    }
    $: onlineUsers = $onlineUserStore;
    $: typingStore = $isTypingStore;
    $: groupTypingStore = $groupIsTypingStore;
    export let isGroup;
    export let userID;
    export let chatID;
    export let userName;
    export let isFirstLoad; // Used only for the first 10 messages fetch
    $: if (isGroup === false) isOnline = onlineUsers.includes(userID);
    let earliestMessageID = 0; // Store last message ID to fetch next messages
    let showEmoji = false;
    const emojis = [
        "😀",
        "😂",
        "🤣",
        "😅",
        "😆",
        "😉",
        "😱",
        "💩",
        "👍",
        "👎",
        "🇪🇪",
        "🐑",
        "👋",
        "🌟",
        "🚀",
        "🎉",
        "😎",
        "🔥",
        "🤖",
        "💯",
    ];
    let textInput = "";
    let inputField;

    $: isTyping = typingStore.includes(userID);

    let allowAudio = true;
    let audio = new Audio("typing.mp3");
    audio.volume = 0.01; //1% volume, DO NOT INCREASE
    audio.loop = true;
    $: if (isTyping && allowAudio) {
        console.log("uo");
        audio.play();
    } else {
        audio.pause();
        audio.currentTime = 0;
    }
    // Get last 10 messages if is primary load
    if (earliestMessageID == 0) {
        let date = new Date();

        fetch(`${API_URL}/messages`, {
            credentials: "include",
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                date: date,
                chat_id: parseInt(chatID, 10),
                message_id: 0, // 0 if first load otherwise last msg id
            }),
        })
            .then((response) => {
                if (response.ok) {
                    return response.json();
                }
            })
            .then((messages) => {
                if (!messages) {
                    return;
                }
                messages = messages.reverse();
                // console.log(messages)
                const chatContainer = document.getElementById(
                    "bottomChatContainer",
                );
                const chatBody = chatContainer.querySelector(
                    `div[chatid="${chatID}"]`,
                );
                if (!chatBody) return;
                messages.forEach((message) => {
                    const messageElem = new Message({
                        target: chatBody,
                        props: {
                            fromUser: message.user,
                            fromUsername: message.username,
                            time: message.date,
                            msgID: message.messageID,
                            msgContent: message.content,
                            AvatarPath: message.avatar,
                        },
                    });
                });
                chatBody.addEventListener("wheel", wheelEventHandler);
                earliestMessageID = messages[0].messageID;
                chatBody.scrollTop =
                    chatBody.scrollHeight - chatBody.clientHeight;
            })
            .catch((error) => {
                console.error(error);
            });
    }

    // SCROLLING (MORE MESSAGES)
    function throttle(func, delay) {
        let throttling = false;
        return function (...args) {
            if (!throttling) {
                throttling = true;
                func.apply(this, args);
                setTimeout(() => {
                    throttling = false;
                }, delay);
            }
        };
    }

    const throttledScroll = throttle(function () {
        const chatContainer = document.getElementById("bottomChatContainer");
        const chatBody = chatContainer.querySelector(`div[chatid="${chatID}"]`);
        const currentScrollPos = chatBody.scrollHeight;
        let date = new Date();
        fetch(`${API_URL}/messages`, {
            credentials: "include",
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                date: date,
                chat_id: parseInt(chatID, 10),
                message_id: earliestMessageID,
            }),
        })
            .then((response) => {
                if (response.ok) {
                    return response.json();
                }
            })
            .then((messages) => {
                if (!messages) {
                    // Edge-case for %10 msg.count
                    chatBody.removeEventListener("wheel", wheelEventHandler);
                    return;
                }
                if (messages.length < 10) {
                    //if the length is < 10, we have the last messages so stop scrolling, if no messages it means there is no history
                    chatBody.removeEventListener("wheel", wheelEventHandler);
                }
                let refMessage =
                    chatBody.querySelectorAll(".message-container")[0]; //top message to insert behind
                messages = messages.reverse(); //reverse order of given message so they are the right way.
                // this is because they are taken from db in reversed order and then the individual 10 messages are put back the right way.
                const fragment = document.createDocumentFragment();
                messages.forEach((message) => {
                    const messageElem = new Message({
                        target: fragment,
                        props: {
                            fromUser: message.user,
                            fromUsername: message.username,
                            time: message.date,
                            msgID: message.messageID,
                            msgContent: message.content,
                            AvatarPath: message.avatar,
                        },
                    });
                });
                chatBody.insertBefore(fragment, refMessage);
                date = messages[0].date; //this tracks the offset of which messages to get
                earliestMessageID = messages[0].messageID; // moved the offset to ID based system
                // Restore scroll pos after new messages
                const addedContentHeight =
                    chatBody.scrollHeight - currentScrollPos;
                chatBody.scrollTop += addedContentHeight;
            })
            .catch((error) => {
                console.error(error);
            });
    }, 500);

    function wheelEventHandler(event) {
        const wheelDirection = event.deltaY < 0 ? 1 : 0;
        const chatContainer = document.getElementById("bottomChatContainer");
        const chatBody = chatContainer.querySelector(`div[chatid="${chatID}"]`);
        if (wheelDirection === 1 && chatBody.scrollTop === 0) {
            throttledScroll();
        }
    }

    // Scrolls to bottom with/without animation
    function scrollToBottom(bodyElem, animation = true) {
        let startTime;
        let start = bodyElem.scrollTop;
        let end = bodyElem.scrollHeight - bodyElem.clientHeight;
        if (animation === false) {
            bodyElem.scrollTop = bodyElem.scrollHeight - bodyElem.clientHeight;
            return;
        }
        let duration = 250;
        function animateScroll(timestamp) {
            if (!startTime) startTime = timestamp;
            const elapsed = timestamp - startTime;
            const progress = Math.min(elapsed / duration, 1);
            bodyElem.scrollTop = start + (end - start) * progress;
            if (elapsed < duration) {
                requestAnimationFrame(animateScroll);
            }
        }
        requestAnimationFrame(animateScroll);
    }

    const throttledTyping = throttle(function () {
        // console.log("SENDING TYPING")
        sendMessage(
            JSON.stringify({
                type: "typing",
                targetid: userID,
                fromid: $userInfo.id,
            }),
        );
    }, 1800);
    const throttledGroupTyping = throttle(function () {
        sendMessage(
            JSON.stringify({
                type: "groupTyping",
                targetid: chatID,
                fromid: $userInfo.id,
                username: $userInfo.firstName + " " + $userInfo.lastName,
            }),
        );
    }, 1800);

    // Handle chat SEND (enter)
    function handleKeyPress(event) {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            const chatBox = event.target.closest(".chatBox");
            const textarea = chatBox && chatBox.querySelector("textarea");
            textarea && (textarea.style.height = "50px");
            // If message is not empty
            if (textInput.trim() !== "") {
                // Compile Message Data to Object (Double obj parsing for msgobj)
                if (isGroup) {
                    let msgObj = JSON.stringify({
                        fromUserID: $userInfo.id,
                        chatID: chatID,
                        content: textInput,
                        AvatarPath: $userInfo.avatar,
                    });

                    sendMessage(
                        JSON.stringify({
                            type: "newMessage",
                            data: msgObj,
                            isgroup: isGroup,
                        }),
                    );
                    markGroupMessageAsSeen(chatID);
                } else {
                    let msgObj = JSON.stringify({
                        fromUserID: $userInfo.id,
                        fromUsername:
                            $userInfo.firstName + " " + $userInfo.lastName,
                        toUserID: userID,
                        chatID: chatID,
                        content: textInput,
                        AvatarPath: $userInfo.avatar,
                    });
                    console.log("Compiled message to send:", msgObj);
                    sendMessage(
                        JSON.stringify({
                            type: "newMessage",
                            data: msgObj,
                            isgroup: isGroup,
                        }),
                    );
                }

                // Scroll chat to bottom after enter is pressed (delay for the message to loop back from backend)
                const chatContainer = document.getElementById(
                    "bottomChatContainer",
                );
                const chatBody = chatContainer.querySelector(
                    `div[chatid="${chatID}"]`,
                );
                setTimeout(() => {
                    scrollToBottom(chatBody);
                }, 160);

                textInput = "";
                event.target.textContent = "";
            }
        } else {
            if (isGroup) {
                throttledGroupTyping();
                return;
            } // Disable IsTyping temporarily for group chats ***TODO
            throttledTyping();
        }
    }

    function emojiBool() {
        showEmoji = !showEmoji;
    }

    function emojiInsert(emoji) {
        textInput += emoji;
    }

    function minimizeChat(event) {
        console.log("test:", event);
        const chatPopup = event.currentTarget.closest(".chat-popup");
        const chatPreview = chatPopup.nextElementSibling;
        chatPopup.classList.remove("chat-popup-open");
        chatPopup.classList.add("chat-popup-close");
        setTimeout(() => {
            chatPopup.style.display = "none";
            chatPopup.classList.remove("chat-popup-close");
            chatPopup.classList.add("chat-popup-open");
        }, 310);
        chatPreview.style.visibility = "visible";
    }

    function toggleChat(event) {
        const chatPreview = event.currentTarget.closest(".chat-preview");
        chatPreview.style.visibility = "collapse";
        const chatPopup = chatPreview.previousElementSibling;
        chatPopup.style.display = "flex";
        // Notification remove on-click
        chatPreview.classList.remove("notification");
        // Check if its the first load (SCROLLING)
        const activeChat = event.currentTarget.closest(".chatBox");
        if (activeChat && activeChat.hasAttribute("isfirstload")) {
            activeChat.removeAttribute("isfirstload");
            const chatBody = chatPopup.querySelector(".chat-body");
            chatBody.scrollTop = chatBody.scrollHeight - chatBody.clientHeight;
        }
    }
    function scrollChatBottom(event) {
        const activeChat = event.currentTarget.closest(".chatBox");
        const chatBody = activeChat.querySelector(".chat-body");
        scrollToBottom(chatBody);
        const notification = activeChat.querySelector(
            ".new-message-notification",
        );
        notification.style.display = "none";

        if (isGroup) {
            markGroupMessageAsSeen(chatID);
        }
    }
    // import svg elements
    let allUsersMap = new Map();
    let users = [];
    let user;
    let chatAvailable = false;
    $: {
        users = $allUsers;
        allUsersMap = new Map(users.map((obj) => [obj.ID, obj]));
        user = allUsersMap.get(userID);
        // Check atleast one of the users follows the other OR if is group
        if (!isGroup) {
            chatAvailable = user
                ? user.IsFollowing === 1 || user.AreFollowing === 1
                : false;
        } else {
            chatAvailable = true;
        }
    }

    function updateHeight() {
        const textarea = this.closest(".chatBox").querySelector("textarea");
        textarea.style.height = "30px";
        textarea.style.height = `${textarea.scrollHeight}px`;
    }
</script>

<div
    class="chatBox"
    {userID}
    {isFirstLoad}
    id="activeChat-chatModule"
    style="display: flex;"
>
    <div class="chat-popup chat-popup-open">
        <div class="chat-header">
            <div class="wrapper">
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                {#if isGroup}
                    <div
                        class="avatar {isOnline ? 'online' : 'offline'}"
                        on:click={() => console.log("suema")}
                    >
                        <img
                            src="{IMAGE_URL}{AvatarPath}"
                            alt={userID}
                            class={isOnline ? "" : "avatar-grayscale"}
                        />
                    </div>
                {:else}
                    <div
                        class="avatar {isOnline ? 'online' : 'offline'}"
                        on:click={() => selectUser(userID)}
                    >
                        <img
                            src="{IMAGE_URL}{AvatarPath}"
                            alt={userID}
                            class={isOnline ? "" : "avatar-grayscale"}
                        />
                    </div>
                {/if}

                <div class="username">
                    <!-- svelte-ignore a11y-missing-attribute -->
                    <a>{userName}</a>
                </div>
            </div>
            <div class="btn-wrapper">
                <!-- Hide/Minimize current chat -->
                <div class="minimize-chat" on:click={(e) => minimizeChat(e)}>
                    <MinimizeChat />
                </div>
                <!-- Close/Remove current chat -->
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <div
                    class="close-chat"
                    on:click={(e) => {
                        removeFromActiveChat(e, "instant", userID);
                        allowAudio = false;
                    }}
                >
                    <CloseChat />
                </div>
            </div>
        </div>
        {#if chatAvailable}
            <!-- messageCount not used REDUNDANCY? ***TODO REMOVE -->
            <div class="chat-body" {chatID} {earliestMessageID} messageCount="">
                <IsTyping {isTyping} {userName} {isGroup} {chatID} />
            </div>
        {:else}
            <ChatFollowing {userID} {userName} {user} />
        {/if}

        <div class="chat-footer">
            {#if chatAvailable}
                <textarea
                    contenteditable
                    class="chatModule-input-field"
                    style="resize: none;"
                    bind:this={inputField}
                    on:input={updateHeight}
                    on:keypress={handleKeyPress}
                    bind:value={textInput}
                />
            {:else}
                <input readonly class="chatModule-input-field" />
            {/if}
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="chatModule-emoji-picker">
                <div on:click={() => emojiBool()}>
                    <ChatModuleEmojiPicker />
                </div>
                {#if showEmoji}
                    <div class="emojiWindow">
                        {#each emojis as emoji}
                            <button
                                on:click={(event) => {
                                    emojiInsert(emoji);
                                    inputField.focus();
                                }}>{emoji}</button
                            >
                        {/each}
                    </div>
                {/if}
            </div>
        </div>
        <div
            class="new-message-notification2 typingGlow"
            style="display:none"
        ></div>
        <div class="new-message-notification" style="display: none;">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div class="notif-wrapper" on:click={scrollChatBottom}>
                <svg
                    width="31px"
                    height="31px"
                    viewBox="0 0 24 24"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    stroke="#ffffff"
                    style="--darkreader-inline-stroke: #e8e6e3;"
                    data-darkreader-inline-stroke=""
                >
                    <g id="SVGRepo_bgCarrier" stroke-width="0"></g>
                    <g
                        id="SVGRepo_tracerCarrier"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    ></g>
                    <g id="SVGRepo_iconCarrier">
                        <path
                            d="M9 13L12 16M12 16L15 13M12 16V8M21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12Z"
                            stroke="#ffffff"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            style="--darkreader-inline-stroke: #e8e6e3;"
                            data-darkreader-inline-stroke=""
                        ></path>
                    </g>
                </svg>
            </div>
        </div>
    </div>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="chat-preview" on:click={toggleChat}>
        <div class="wrapper">
            <div class="avatar {isOnline ? 'online' : 'offline'}">
                <img
                    src="{IMAGE_URL}{AvatarPath}"
                    alt={userID}
                    class={isOnline ? "" : "avatar-grayscale"}
                />
            </div>
            <div class="username">
                <a id="preview-username">{userName}</a>
            </div>
        </div>
        <div class="btn-wrapper">
            <div class="close-chat">
                <CloseChat />
            </div>
        </div>
    </div>
    <script>
    </script>
</div>

<style>
    /* ajutine: */
    #preview-username,
    .username a {
        color: #ffffff;
        text-decoration: none;
        font-size: 16px;
        margin-right: 10px;
        margin-left: 30px;
    }
    :root {
        --chatWidth: 264px;
        --chatPreviewH: 40px;
        --chatFullH: 400px;
    }

    .chatBox {
        margin-right: 6px;
        margin-left: 6px;
        width: var(--chatWidth);
        height: (--chatPreviewH);
        display: flex;
        flex-direction: row;
        margin-right: 5px;
        /* is needed because of the emojiWindow */
    }

    .chat-preview {
        display: flex;
        visibility: collapse;
        justify-content: space-between;
        width: var(--chatWidth);
        height: var(--chatPreviewH);
        background-color: #011;
        border-radius: 5px;
        border: 1px solid rgb(145, 145, 145);
    }
    .chat-preview .wrapper {
        width: fit-content;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .chat-preview .btn-wrapper {
        height: 100%;
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .chat-popup {
        display: flex;
        flex-direction: column;
        position: absolute;
        transform: translatey(-350px);
        width: var(--chatWidth);
        height: var(--chatFullH);
        background-color: #011;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.125);
    }
    .chat-popup-open {
        animation: chat-open 300ms ease-in-out forwards;
        transform-origin: bottom;
    }

    :global(.chat-popup-close) {
        animation: chat-close 300ms ease-in-out forwards;
        transform-origin: bottom;
    }
    .chatModule-emoji-picker {
        margin-right: 3px;
        width: 34px;
        height: 34px;
    }

    .emojiWindow {
        position: absolute;
        display: flex;
        flex-wrap: wrap;
        justify-content: flex-start;
        padding: 10px;
        background-color: #011;
        border-radius: 12px;
        width: 50%;
        max-width: 600px;
        border: 1px solid rgba(255, 255, 255, 0.125);
        bottom: 100%;
        margin-bottom: 10px;
        right: 15px; /* Align to the right edge of .chatBox */
    }

    .emojiWindow button {
        flex: 0 0 calc(10% - 2px);
        margin: 5px;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 15px;
        box-sizing: border-box;
    }
    .avatar {
        cursor: pointer;
        margin-left: 6px;
        margin-top: 3px;
        margin-bottom: 3px;
        width: 25px;
        height: 25px;
        border-radius: 50%;
        /* border: 2px solid green; */
    }

    .avatar img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 50%;
    }
    .chat-header {
        display: flex;
        justify-content: space-between;
        width: var(--chatWidth);
        height: 60px;
        background-color: #011;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.125);
        border-bottom-left-radius: 1px;
        border-bottom-right-radius: 1px;
    }

    .chat-header .wrapper {
        width: fit-content;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-left: 12px;
    }

    .chat-header .btn-wrapper {
        height: 100%;
        position: relative;
        right: 0;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .chat-header .avatar {
        margin-left: 0;
    }

    .chat-body {
        width: var(--chatWidth);
        height: 280px;
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        overflow: hidden scroll;
    }

    /* :global(.message-container) {
        margin: 2px;
        display: flex;
        align-content: left;
        margin-top: 2px;
        margin-bottom: 2px;
        width: var(--chatWidth);
        min-height: var(--chatPreviewH);
        height: fit-content;
    } */
    /* :global(.message-header) {
        display: flex;
        justify-content: space-between;
    } */

    .chat-footer {
        position: absolute;
        bottom: 0;
        left: 0;
        width: var(--chatWidth);
        height: var(--chatPreviewH);
        background-color: #011;
        border-radius: 12px;
        border: 1px solid rgba(255, 255, 255, 0.125);
        display: flex;
        flex-direction: row;
        align-items: center;
        width: 100%;
        min-height: var(--chatPreviewH);
        height: fit-content;
        border-top-left-radius: 0;
        border-top-right-radius: 0;
        z-index: 2;
    }
    .chatModule-input-field {
        overflow: hidden;
        height: 30px;
        margin: 5px;
        backdrop-filter: blur(16px) saturate(180%);
        -webkit-backdrop-filter: blur(16px) saturate(180%);
        background-color: #011;
        border-radius: 12px;
        border: 1px solid black;
        width: 220px;
        text-align: left;
    }
    .chatModule-input-field[contenteditable] {
        color: white;
        font-weight: 1000;
        min-height: 18px;
        padding: 8px;
        padding-top: 4px;
        padding-bottom: 6px;
        line-height: 1.4;
        border-radius: 10px;
    }
    .chatModule-input-field[contenteditable]:focus {
        outline: none;
    }

    @keyframes user-active-chat-remove {
        0% {
            opacity: 1;
            transform: scaleY(1);
            transform: translatey(0px);
        }
        100% {
            opacity: 0;
            transform: scaleY(0);
            transform: translatey(-20px);
        }
    }
    :global(.user-active-chat-remove) {
        animation: user-active-chat-remove 300ms ease-in-out forwards;
        transform-origin: top;
    }
    @keyframes chat-open {
        0% {
            opacity: 0;
            height: var(--chatPreviewH);
            transform: scaleY(0);
            transform: translatey(0px);
        }
        100% {
            opacity: 1;
            height: var(--chatFullH);
            transform: scaleY(1)
                translateY(calc(-1 * (var(--chatFullH) - var(--chatPreviewH))));
        }
    }

    @keyframes chat-close {
        0% {
            opacity: 1;
            height: var(--chatFullH);
            transform: scaleY(1)
                translateY(calc(-1 * (var(--chatFullH) - var(--chatPreviewH))));
        }
        100% {
            opacity: 0;
            height: var(--chatPreviewH);
            transform: scaleY(0);
            transform: translatey(0px);
        }
    }

    :global(.full-name) {
        position: absolute;
        transform: translate(20px, 14px);
        font-family: "Jura";
        color: white;
        font-size: small;
        font-weight: 1000;
    }
    @keyframes wave {
        0% {
            transform: translateY(-0.5px);
        }
        45% {
            transform: translateY(-4px);
        }
        100% {
            transform: translateY(0px);
        }
    }
    /*     :global(.typingGlow) {
		display: block;
		background: linear-gradient(0deg, rgba(150, 4, 254, 0.645) 43%, rgba(178,4,254,0) 92%); 
		animation: pulseGlow 1.5s infinite;
	} */
    @keyframes pulseGlow {
        0% {
            opacity: 0.2;
        }
        50% {
            opacity: 1;
        }
        100% {
            opacity: 0;
        }
    }

    /* ||> Notifications */
    .new-message-notification {
        display: none;
        width: 100%;
        height: 99px;
        background: rgb(178, 4, 254);
        background: linear-gradient(
            0deg,
            rgba(178, 4, 254, 0.7805322812718838) 43%,
            rgba(178, 4, 254, 0) 92%
        );
        position: absolute;
        bottom: 0px;
        z-index: 1;
        border-bottom-left-radius: 12px;
        border-bottom-right-radius: 12px;
        animation: pulseBorder 0.22s cubic-bezier(0.23, 1, 0.32, 1);
    }
    .new-message-notification2 {
        display: none;
        width: 100%;
        height: 99px;
        background: rgb(178, 4, 254);
        background: linear-gradient(
            0deg,
            rgba(178, 4, 254, 0.7805322812718838) 43%,
            rgba(178, 4, 254, 0) 92%
        );
        position: absolute;
        bottom: 0px;
        z-index: 1;
        border-bottom-left-radius: 12px;
        border-bottom-right-radius: 12px;
        animation: pulseBorder 0.22s cubic-bezier(0.23, 1, 0.32, 1);
    }

    @keyframes pulseBorder {
        0% {
            background: linear-gradient(
                0deg,
                rgba(178, 4, 254, 0.7805322812718838) 22%,
                rgba(178, 4, 254, 0) 92%
            );
        }
        25% {
            background: linear-gradient(
                0deg,
                rgba(178, 4, 254, 0.7805322812718838) 43%,
                rgba(178, 4, 254, 0) 92%
            );
        }
        50% {
            background: linear-gradient(
                0deg,
                rgba(178, 4, 254, 0.7805322812718838) 66%,
                rgba(178, 4, 254, 0) 92%
            );
        }
        75% {
            background: linear-gradient(
                0deg,
                rgba(178, 4, 254, 0.7805322812718838) 52%,
                rgba(178, 4, 254, 0) 92%
            );
        }
        100% {
            background: linear-gradient(
                0deg,
                rgba(178, 4, 254, 0.7805322812718838) 43%,
                rgba(178, 4, 254, 0) 92%
            );
        }
    }

    .notif-wrapper {
        cursor: pointer;
        position: absolute;
        top: 4px;
        right: 6px;
        width: 35px;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    :global(.notification) {
        animation:
            pulse 1s forwards,
            shake3 0.6s forwards,
            glow 1s forwards;
    }

    @keyframes shake3 {
        0%,
        100% {
            transform: translateX(0);
        }

        10%,
        30%,
        50%,
        70%,
        90% {
            transform: translateX(-3px);
        }

        20%,
        40%,
        60%,
        80% {
            transform: translateX(3px);
        }
    }

    @keyframes pulse {
        0% {
            background-color: rgba(17, 25, 40, 0.403);
        }
        50% {
            background-color: rgba(58, 71, 95, 0.745);
        }
        100% {
            background-color: rgba(32, 40, 56, 0.745);
        }
    }

    @keyframes glow {
        100% {
            box-shadow: 0px 0px 10px 1px rgba(201, 16, 230, 0.75);
        }
    }
</style>
