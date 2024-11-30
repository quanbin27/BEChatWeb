class Store {
    constructor() {
    }

    setChat(chat){
        this.chat = chat;
    }
    getContacts(){
        chat.getContacts();
    }
    setContact(contact){
        this.contact =contact;
    }
    getContact(){
        return this.contact;
    }
    setGroup(group){
        this.group = group;
    }
    getRroup(){
        return this.group;
    }
    setProfile(profile){
        this.profile = profile;
    }
    getProfile(){
        return this.profile;
    }
}


class Chat{
    constructor(store){
        this.store = store;

        this.userConversation = document.querySelector('#users-conversation');
        this.btnSend = document.querySelector('.links-list-item button[type="submit"]');
        this.chatinput = document.querySelector('.chat-input');

        this.favouriteUsers= document.querySelector('#favourite-users');
        this.favouriteGroups = document.querySelector('#favourite-groups');

        this.userName = document.querySelector('.user-profile-show');
        this.avatar = document.querySelector('.avatar-sm');
        this.profileDetail = document.querySelector('.user-profile-sidebar');
        this.profileAvatar = this.profileDetail.querySelector('.profile-img');
        this.profileUsername = this.profileDetail.querySelector('.user-name.mb-1.text-truncate');
        this.profileUsername1 = this.profileDetail.querySelector('.user-name.font-size-14.text-truncate');

        this.current_conversation_id = null;
    }

    async init(){
        this.favouriteUserData = await this.getFavouriteUsers();
        this.favouriteGroupData = await this.getFavouriteGroups();

        this.renderFavouriteUsers(this.favouriteUserData);
        this.renderFavouriteGroups(this.favouriteGroupData);


        if(this.favouriteUserData.length>0){
            this.selectedUser(this.favouriteUserData[0])
        }
        else if(this.favouriteGroupData.length>0){
            this.selectedGroup(this.favouriteGroupData[0]);
        }
        else{
            this.renderDefaultMessage();
        }
            
        this.btnSend.onclick = async (event)=>{
            event.preventDefault();
            if(this.chatinput.value == '')
                return;
            let newMessage = {
                group_id: this.current_conversation_id,
                content: this.chatinput.value, 
            }
            if(this.current_conversation_id === null) return;
            this.chatinput.value = '';
            const response = await this.sendMessage(newMessage);
            if(response === 'error'){
                return;
            }
            newMessage = {
                ID: response.MessageID,
                UserID: this.store.profile.user_id,
                GroupID: this.current_conversation_id,
                Content: newMessage.content,
            }
            this.receiveMessage(newMessage.GroupID,newMessage.UserID,newMessage.Content,newMessage.ID);
            
        }
    }

    
    renderDefaultMessage() {
        const defaultMessages = [
            { Content: 'Welcome to our website!' },
            { Content: 'Feel free to explore and chat.' },
            { Content: 'You can add friends and create groups.' }
        ];
        const btnAddMoreContact = document.querySelector('.add-more-contact');
        btnAddMoreContact.style.display = 'none';
        let image_path = 'assets/images/users/user-dummy-img.jpg';
        this.avatar.src = image_path;
        this.userName.textContent = 'Default chat';
        this.profileUsername.textContent = 'Default chat';
        this.profileUsername1.textContent = 'Default chat';
        this.profileAvatar.src = image_path;

        const groupContact = this.profileDetail.querySelector('.group-contacts');
        groupContact.style.display = 'none';

        this.renderConversation(defaultMessages);
    }
    
    async sendMessage(messageData){
        let token = this.store.profile.getToken();

        return fetch('/api/v1/message',{
            method:'POST',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            },
            body:JSON.stringify(messageData)
        })
        .then(async response => {
            if(response.ok){
                const responseJson = await response.json();
                return responseJson;
            }
            return 'error';
        })
        .catch(()=>{
            return 'error';
        })

    }
    async getCurrentConversationById(conversationId){
    
        let token = this.store.profile.getToken();

        const conversation = fetch(`/api/v1/group/${conversationId}/message`,{
            method:'get',
            headers:{
                'authorization':`Bearer ${token}`
            }
        })
        .then(async response=>{
            if(response.ok){
                const conversationData = await response.json();
                if(Object.keys(conversationData).length == 0) return [];
                return conversationData.Messages;
            }
            return []
        }).catch(()=>{return []})
        return conversation;
    }
    async getFavouriteUsers(){
        // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        let token = this.store.profile.getToken();

        const users = fetch('/api/v1/user/group-chat/2',{
            method:'get',
            headers:{
                'authorization':`Bearer ${token}`
            }
        })
        .then(async response=>{
            if(response.ok){
                const groups = await response.json();
                if(Object.keys(groups).length == 0)return [];
                // let fullGroupData = []
                for(let gr of groups.groups){
                    const otherUserId = gr.other_user_id;
                    const otherUserInfo = await this.store.contact.getContactInfo(otherUserId);
                    gr.username = otherUserInfo.Name;
                    gr.image_path = otherUserInfo.avatar;
                }
                return groups.groups;
            }
            return []
        }).catch(()=>{return []})
        return users;
    }
    async getFavouriteGroups(){
        // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        let token = this.store.profile.getToken();

        const groups = fetch('/api/v1/user/group-chat/3',{
            method:'get',
            headers:{
                'authorization':`Bearer ${token}`
            }
        })
        .then(async response=>{
            if(response.ok){
                const groups = await response.json();
                if(Object.keys(groups).length == 0 )return [];
               
                return groups.groups;
            }
            return []
        }).catch(()=>{return []})
        return groups;
    }
    messageElement(conversation){
        const chatListRight = document.createElement('li');
        chatListRight.setAttribute('data-message-id',conversation.ID);
        
        const now = new Date();
        const hours = now.getHours(); // Lấy giờ
        const minutes = now.getMinutes(); // Lấy phút
        const seconds = now.getSeconds(); // Lấy giây
        chatListRight.innerHTML =`
        <div class="conversation-list">
            <div class="user-chat-content">
                ${conversation.UserID===this.store.profile.user_id? '':'<span>${Võ văn a}</span>'}
                
                <div class="ctext-wrap">
                    <div class="ctext-wrap-content">
                        <p class="mb-0 ctext-content">${conversation.Content}</p>
                    </div>
                    <div class="dropdown align-self-start message-box-drop">
                        <a class="dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            <i class="ri-more-2-fill"></i>
                        </a>
                        <div class="dropdown-menu">
                           
                            <a class="dropdown-item d-flex align-items-center justify-content-between delete-item" id="delete-item-1" href="#">
                                Delete <i class="bx bx-trash text-muted ms-2"></i>
                            </a>
                        </div>
                    </div>
                </div>
                <div class="conversation-name">
                    <small class="text-muted time">${hours}:${minutes} pm</small>
                    <span class="text-success check-message-icon">
                        <i class="bx bx-check"></i>
                    </span>
                </div>
            </div>
        </div>
        `
        if(conversation.UserID == this.store.profile.user_id){
            chatListRight.className = 'chat-list right'
            const btnDelete = chatListRight.querySelector('.delete-item');
                btnDelete.onclick = async ()=>{
                    const messageResponse = await this.deleteMessage(messageData.ID);
                    if(messageResponse==='success'){
                        this.userConversation.removeChild(message);
                    }
                }
        }
        else{
            chatListRight.className = 'chat-list left';
        }
        return chatListRight;
    }

    // dùng để chọn chat bên contact;
    renderConversationOfContact(state){

        const userId = state.id;
        let userData = this.favouriteUserData.find(user=>user.userId == userId);
        if(userData==null){
            userData = {
                id:null,//tạo sau
                userId:state.id,
                name:state.name,
                email:state.email,
                image_path:state.image_path,
            }}
        this.selectedUser(userData);
    }
    async deleteMessage(messageId){
        // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        let token = this.store.profile.getToken();
        return fetch('/api/v1/message',{
            method:'DELETE',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            },
            body:JSON.stringify({
                "message_id": messageId
            })
        })
        .then( async response =>{
            if(response.ok){
                const message = await response.json();
                return 'success';
            }
            else return 'error';
        }).catch(()=>{return 'error';})
    }
    renderConversation(current_conversation){
        console.log(current_conversation);
        this.userConversation.innerHTML = ''
        for(let messageData of current_conversation){
            const message = this.messageElement(messageData);
            this.userConversation.appendChild(message);
        }
        const scroll = document.querySelector('.chat-conversation .simplebar-offset .simplebar-content-wrapper')
        scroll.scrollTop = scroll.scrollHeight;
    }
    decodeImage(){
        // tạm
    }

    findFavouriteUserAndGroupLi(chatId){
        for(let li of this.favouriteUsers.querySelectorAll('li')){
            let group_id = li.getAttribute('data-chat-id');
            if(parseInt(group_id) === parseInt(chatId)){
                return {
                    'li': li,
                    'type':'user',
                }
            }
        }
        
        for(let li of this.favouriteGroups.querySelectorAll('li')){
            let group_id = li.getAttribute('data-chat-id');
            if(parseInt(group_id) === parseInt(chatId)){
                return {
                    'li': li,
                    'type':'group',
                }
            }
        }
        return null;    
    }
    createFavouriteUserLi(userData){
        
        const user = document.createElement('li');
        user.style.cursor = 'pointer';
        user.onmouseover=()=>{
            user.style.backgroundColor = '#d5d5d5'
        }
        user.onmouseleave=()=>{
            user.style.backgroundColor = 'white';
        }
        user.style.marginTop = '10px';
        if(userData.image_path == null){
            userData['image_path'] = 'assets/images/users/user-dummy-img.jpg'
        }
        if(userData.message == null) userData['message'] = '';
        user.innerHTML = `
            <div style='display:flex; gap:10px; margin-left:10px; align-items:center  '>
                <img src = ${userData.image_path} style = 'width:2.4rem; height:2.4rem; border-radius:50%!important;'>
                <div style = 'display:flex; flex-direction:column'>
                    <span style= 'font-size: medium; font-weight:600;' class = 'username'>${userData.username}</span>
                    <span class = 'last-message' style = 'text-overflow: ellipsis; white-space: nowrap; overflow: hidden; max-width:190px;' >${userData.latest_message==null?'':userData.latest_message}</span>
                </div>
            </div>
        `
        user.setAttribute('data-user-id',userData.other_user_id);//chua co
        user.setAttribute('data-chat-id',userData.group_id);
        user.onclick = ()=>{
            this.selectedUser(userData)
        }
        return user;
    }
    async renderFavouriteUsers(favouriteUserData){

        this.favouriteUsers.innerHTML = '';
        for(let userData of favouriteUserData){
            const user = this.createFavouriteUserLi(userData);
            this.favouriteUsers.appendChild(user);

        }
    }

    async selectedUser(userData){
        let user = this.favouriteUsers.querySelector(`li[data-chat-id="${userData.group_id}"]`);
        if(user == null)
            user = this.favouriteGroups.querySelector(`li[data-chat-id="${userData.group_id}"]`);
        if(user == null) alert('lỗi ở select chat');
        const btnAddMoreContact = document.querySelector('.add-more-contact');
        btnAddMoreContact.style.display = 'none';

        this.avatar.src = userData.image_path;
        this.userName.textContent = userData.username;
        this.profileUsername.textContent = userData.username;
        this.profileUsername1.textContent = userData.username;
        this.profileAvatar.src = userData.image_path;

        const groupContact = this.profileDetail.querySelector('.group-contacts');
        groupContact.style.display = 'none';

        const currentConversationId = userData.group_id;
        this.current_conversation_id = currentConversationId;
        const current_conversation = await this.getCurrentConversationById(currentConversationId);
        this.renderConversation(current_conversation);
        let currentChat = this.favouriteUsers.querySelector('li.current');
        if(currentChat === null) currentChat = this.favouriteGroups.querySelector('li.current');
        if(currentChat !=null)
        {
            currentChat.style.backgroundColor = 'white';
            currentChat.classList.remove('current');
            currentChat.onmouseleave=()=>{
                currentChat.style.backgroundColor = 'white';
            }
        }
        user.style.backgroundColor = '#d5d5d5';
        user.classList.add('current');
        user.onmouseleave=()=>{

        }
    }

    createFavouriteGroupLi(groupData){
        const group = document.createElement('li');
        group.style.cursor = 'pointer';
        group.onmouseover=()=>{
            group.style.backgroundColor = '#d5d5d5'
        }
        group.onmouseleave=()=>{
            group.style.backgroundColor = 'white';
        }
        group.style.marginTop = '10px';
        let image_path = 'assets/images/group.jpg';
        group.innerHTML = `
            <div style='display:flex; gap:10px; margin-left:10px; align-items:center  '>
                <img src = ${image_path} style = 'width:2.4rem; height:2.4rem; border-radius:50%!important;'>
                <div style = 'display:flex; flex-direction:column'>
                    <span style= 'font-size: medium; font-weight:600;' class = 'groupname'>${groupData.group_name}</span>
                    <span  class = 'last-message' style = 'text-overflow: ellipsis; white-space: nowrap; overflow: hidden; max-width:180px;'>${groupData.latest_message}</span>
                </div>
            </div>
        `
        group.setAttribute('data-chat-id',groupData.group_id);
        group.onclick = ()=>{
            this.selectedGroup(groupData)
        }
        return group;
    }    
    renderFavouriteGroups(favouriteGroupData){
        this.favouriteGroups.innerHTML = '';
        let image_path = 'assets/images/group.jpg';
        for(let groupData of favouriteGroupData){
            groupData['image_path'] = image_path;
            const group = this.createFavouriteGroupLi(groupData);
            this.favouriteGroups.appendChild(group);

        }
    }

    async renderMemberInGroup(currentConversationId){
        function createLi(member){
            let li = document.createElement('li');
            li.style.display = 'flex';
            li.style.alignItems='center';
            li.style.justifyContent = 'space-between';
            li.innerHTML = `
                <div>
                    <p style = 'margin:0'>${member.Name}</p>
                    <p style = 'margin:0'>${member.RoleID ==1?'admin':''}</p>
                </div>
                <div class="d-flex">
                    <div class="flex-shrink-0">
                        <div class="dropdown">
                            <button class="btn nav-btn text-black dropdown-toggle" type="button" data-bs-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                                <i class='bx bx-dots-vertical-rounded'></i>
                            </button>
                            <div class="dropdown-menu dropdown-menu-end">
                                <a class="btn-delete dropdown-item d-flex justify-content-between align-items-center" href="#">Delete <i class="bx bx-trash text-muted"></i></a>
                                <a class="btn-role dropdown-item d-flex justify-content-between align-items-center" href="#">Switch admin<i class="bx bxs-pencil text-muted"></i></a>
                            </div>
                        </div>
                    </div>
                </div>
            `
            return li;
        }
        async function deleteMembers(groupId,id,token){
            return fetch(`/api/v1/group/${groupId}/member`,{
                method:'DELETE',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body: JSON.stringify({
                    member_id:id
                })
            }).then(async re=>{
                if(re.ok){
                    return 'success';
                }
                return 'error';
            })
            .catch(()=>{
                return 'error';
            })
        }
        // async leaveGroup(groupId)//tạm
        if(currentConversationId == null){
            currentConversationId = this.current_conversation_id;
        }
        const btnAddMoreContact = document.querySelector('.add-more-contact');
        
        const groupContact = this.profileDetail.querySelector('.group-contacts');
        groupContact.style.display = '';
        const ul = this.profileDetail.querySelector('.list-group-members');
        ul.innerHTML = ''

        const members = await this.store.contact.getContactInGroup(currentConversationId);
        let role = 'member';
        const admin = members.find(member=>member.RoleID ==1)
        if(admin.UserID == this.store.profile.user_id){
            role = 'admin';
            btnAddMoreContact.style.display = '';
        } 
        
        for(let member of members){
            let li = createLi(member);
            const btnWrapper = li.querySelector('.d-flex');
            const temp = btnWrapper.querySelector('.dropdown-menu.dropdown-menu-end')
            const btnRole = temp.querySelector('.btn-role');
            const btnDelete = temp.querySelector('.btn-delete');
            if(role === 'member'){
                if(member.UserID==this.store.profile.user_id){
                    temp.removeChild(btnRole);
                    btnDelete.onclick = ()=>{
                        alert(currentConversationId);
                    }
                }
                else
                    li.removeChild(btnWrapper);
            }
            else{
                btnDelete.onclick = async ()=>{
                    let token = this.store.profile.getToken();
                    const message = await deleteMembers(currentConversationId,member.UserID,token);
                    if(message === 'success'){
                        ul.removeChild(li);
                        toastr.success('Delete successfully','Success')
                    }else{
                        toastr.error('Can not delete this member now.','Error')
                    }
                }
                if(member.UserID==this.store.profile.user_id){
                    li.removeChild(btnWrapper);
                }
            }
            ul.appendChild(li);
        }
    }
    async selectedGroup(chatData){


        let chat = this.favouriteUsers.querySelector(`li[data-chat-id="${chatData.group_id}"]`);
        if(chat == null)
            chat = this.favouriteGroups.querySelector(`li[data-chat-id="${chatData.group_id}"]`);
        if(chat == null) alert('lỗi ở select chat');
        
        this.userName.textContent = chatData.group_name;
        this.profileUsername.textContent = chatData.group_name;
        this.profileUsername1.textContent = chatData.group_name;
        this.profileAvatar.src = chatData.image_path;
        this.avatar.src = chatData.image_path;

        const currentConversationId = chatData.group_id;

        this.current_conversation_id = currentConversationId;
   
        this.renderMemberInGroup(currentConversationId);
        const current_conversation = await this.getCurrentConversationById(currentConversationId);
        this.renderConversation(current_conversation);
        let currentChat = this.favouriteUsers.querySelector('li.current');
        if(currentChat === null) currentChat = this.favouriteGroups.querySelector('li.current');
        if(currentChat !=null)
        {
            currentChat.style.backgroundColor = 'white';
            currentChat.classList.remove('current');
            currentChat.onmouseleave=()=>{
                currentChat.style.backgroundColor = 'white';
            }
        }
        chat.style.backgroundColor = '#d5d5d5';
        chat.classList.add('current');
        chat.onmouseleave=()=>{

        }
    }

    findFavouriteUserLiByUserId(userId){
        const list = this.favouriteUsers.querySelectorAll('li');
        for(const li of list){
            const id = li.getAttribute('data-user-id');
            if(parseInt(id) === parseInt(userId)){
                return li;
            }
        }
        return null;
    }
    findFavouriteGroupLiByGroupId(groupId){
        const list = this.favouriteGroups.querySelectorAll('li');
        for(const li of list){
            const id = li.getAttribute('data-chat-id');
            if(parseInt(id) === parseInt(groupId)){
                return li;
            }
        }
        return null;
    }
    //thêm vào để xử lý khi có sự kiện thêm liên hệ từ socket nếu liên hệ đó có rồi thi thôi không thì thêm
    async addFavoriteUsers(userId){
        console.log('select new contact',userId);
        let li = this.findFavouriteUserLiByUserId(userId);
        if(li!== null){
            return;
        }
        
        const favouriteUsers = await this.getFavouriteUsers();
        this.renderFavouriteUsers(favouriteUsers);
    }
    // thêm group mới vào khi người dùng đó được người khác thêm vào.
    async addNewGroupToFavoriteGroups(groupId,groupName){
        console.log(groupId,groupName);
        const groupData = {
            group_id:groupId,
            group_name:groupName,
            latest_message: 'You have been added to the group abc'
        }
        let li = this.findFavouriteGroupLiByGroupId(groupId);
        if(li!=null) return;
        li = this.createFavouriteGroupLi(groupData);
        this.favouriteGroups.insertBefore(li,this.favouriteGroups.firstChild);
    }

    // xóa group ra khỏi giao diện người dùng khi người dùng đó bị click hoặc tự out
    async deleteGroupInFavoriteGroups(groupId){
        const li = this.findFavouriteGroupLiByGroupId(groupId);
        if(li == null)return;
        this.favouriteGroups.removeChild(li);

        // nhằm xóa đi avatar hay thành viên còn sót lại khi bị out nhóm
        if(this.current_conversation_id == groupId)
            this.renderDefaultMessage();
    }

    findMessageLiInConversation(messageId){
        const list = this.userConversation.querySelectorAll('li');
        for(let li of list){
            let dataMessageId = li.getAttribute('data-message-id');
            if(parseInt(dataMessageId) === parseInt(messageId)){
                return li;
            }
        }
        return null;
    }
    deleteMessageOfOtherConversation(groupId,messageId){
        console.log('groupId:   ',groupId, 'messageId: ',messageId)
        if(this.current_conversation_id !== groupId) return;
        const li = this.findMessageLiInConversation(messageId);
        if(li == null) return;
        li.querySelector('.mb-0.ctext-content').textContent = 'The user has deleted the message.'
        li.querySelector('.ctext-wrap-content').style.backgroundColor = 'white';
        // if(li!=null)
        //     this.userConversation.removeChild(li);
    }
    receiveMessage(groupId,userId,content,messageId){
        console.log(userId)
        const li = this.findFavouriteUserAndGroupLi(groupId);

        if(li == null) return;
        const lastMessage = li.li.querySelector('.last-message');
        lastMessage.textContent = content;
        if(li.type === 'user'){
            if(this.favouriteUsers.firstChild != li.li) {
                this.favouriteUsers.removeChild(li.li);
                this.favouriteUsers.insertBefore(li.li,this.favouriteUsers.firstChild);
            }
        }
        else{
            if(this.favouriteGroups.firstChild != li.li){
                this.favouriteGroups.removeChild(li.li);
                this.favouriteGroups.insertBefore(li.li,this.favouriteGroups.firstChild);
            }
        }
        if(this.store.chat.current_conversation_id != groupId) {
            return;
        }
        const newMessage = this.messageElement({
            ID:messageId,
            UserID: userId,
            Content:content
        })
        this.userConversation.appendChild(newMessage);
        const scroll = document.querySelector('.chat-conversation .simplebar-offset .simplebar-content-wrapper')
        scroll.scrollTop = scroll.scrollHeight;
    }
}

class Group{
    constructor(store){
        this.store  = store;
        this.groupModal = document.querySelector('#addgroup-exampleModal');
        this.addGroupBtn = document.querySelector('#add-group-btn');
        this.btnClose1 = this.groupModal.querySelector('.btn-close');
        this.btnClose2 = this.groupModal.querySelector('.btn-link');
        this.btnCreateGroup = this.groupModal.querySelector('.btn.btn-primary');
        this.addGroupBtn.onclick = ()=>{
            this.openForm();
        }
        this.btnClose1.onclick = this.btnClose2.onclick = ()=>this.closeForm();
    }
    getFriendList(){
        return this.store.getContact().getContacts();
    }
    closeForm(){
        document.body.classList.remove('modal-open');
        document.body.style.overflow = '';
        document.body.removeChild(this.background)
        this.groupModal.style.display = 'none';
        this.groupModal.classList.remove('show');
    }
    async openForm(){
        function createContactElement(contact){
            console.log(contact)
            const li =document.createElement('li');
            li.style.display='flex';
            li.style.gap='10px';
            li.style.alignItems = 'center';
            // li.className = 'form-check';
            li.innerHTML = `
                <input type="checkbox" class="form-check-input">
                <img src = ${contact.avatar} alt = 'not-found' style="width: 2.4rem;height: 2.4rem; border-radius: 50%">
                <p style="margin: 0; padding:0;">${contact.username}</p>
            `
            return li;
        }
        function createGroup(group_info,token){
          
            const response = fetch('/api/v1/group',{
                method:'POST',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body:JSON.stringify(group_info)
            }).then(async re=>{
                console.log(re);
                const message = await re.json()
                
                if(re.ok ){
                    return 'success';
                }
                return 'error';
            }).catch((error)=>{
               
                return 'error';
            })
            return response;
        }
        this.background = document.createElement('div');
        this.background.className = 'modal-backdrop fade show';
        document.body.appendChild(this.background);
        document.body.classList.add('modal-open');
        document.body.style.overflow = 'hidden';
        this.groupModal.style.display = 'block';
        this.groupModal.classList.add('show');

        const contactTag = this.groupModal.querySelector('.card-body.p-2 .simplebar-content .list-unstyled.contact-list');
        contactTag.innerHTML = ''
        const addGroupNameInput = this.groupModal.querySelector('#addgroupname-input');
        let listCheckedFriends = []
        let contacts = await  this.getFriendList();
       
        for( let contact of contacts){
            // Tạo phần tử <li>
            const li = createContactElement(contact);

            // Tạo phần tử <input> với class "form-check-input"
            const input = li.querySelector('input');
            input.onchange = ()=>{
                if(input.checked){
                    listCheckedFriends.push(contact.user_id);
                }
                else
                    listCheckedFriends = listCheckedFriends.filter(id => id !== contact.user_id);
               

            }

            contactTag.appendChild(li); // Gắn <li> vào danh sách <ul>


        }
        this.btnCreateGroup.onclick = async ()=>{
            if(addGroupNameInput.value === ''){
                toastr.info('Please enter group name', 'Notification')
                return;
            }
            if(listCheckedFriends.length<2){
                toastr.info('Please choose at least 2 members','Notification');
                return;
            }
            const groupInfo = {
                name: addGroupNameInput.value,
                member_ids: listCheckedFriends,
            }
            let token = this.store.profile.getToken();
            const message = await createGroup(groupInfo,token);
          
            if(message === 'success'){
                toastr.success('Create group successfully', 'Success');

                const favouriteGroupData = await this.store.chat.getFavouriteGroups();
                this.store.chat.renderFavouriteGroups(favouriteGroupData);
                this.closeForm();
            }
            else{
                toastr.error('Can not create group', 'Failure');

            }
        }
    }

    sendData(){

        //gửi dữ liệu đi;
    }
}

class Contact{
    constructor(store){
        this.store = store;
        
        this.background = document.createElement('div');
        this.background.className = 'modal-backdrop fade show';


        this.btnOpenFormAddContact = document.querySelector('.btn.btn-soft-primary.btn-sm');
        this.addContactModal = document.querySelector('#addContact-exampleModal');
        this.btnSearchContact = this.addContactModal.querySelector('.modal-footer .btn.btn-primary');
        this.findedContactWrapper = this.addContactModal.querySelector('.contact-info');

        this.contactDiv = document.querySelector('.sort-contact');
        this.addMoreContactModal = document.querySelector('.modal.fade.addMoreContactModal')
        this.btnAddMoreContact = document.querySelector('.add-more-contact');
        

        this.pendingContactForm = document.querySelector('#pending-contact-modal');
        this.btnPendingContacts = document.querySelector('#btn-pending-contacts');
        
    }
    //========================contact info ==============================
    async getContactInfo(user_id){
        return fetch(`/api/v1/user/${user_id}`)
        .then(response =>{
            if(response.ok){
                return response.json(); 
            }
            return {}
        })
        .catch(()=>{
            return 'error';
        })
    }
    //=========================== init=======================================
    async init(){
        const contacts = await this.getContacts();

        this.renderContacts(contacts);

        this.btnAddMoreContact.onclick = ()=>{
            this.openFormAddMoreContactModel();
        }

        // gắn sự kiện cho nút xem người đang chờ kết bạn
        this.btnPendingContacts.onclick = ()=>{
            this.openPedingContactForm()
        }

        // gắn sự kiện tìm người liên hệ
        this.btnSearchContact.onclick = ()=>{
            this.findContactByEmail();
        }

    }
    // =======================Pending Contact================================
    getPendingContacts(){
        let token = this.store.profile.getToken();
        const response = fetch('/api/v1/contact/pending-received',{
            method:'GET',
            headers:{
                'authorization':`Bearer ${token}`
            }
        }).then(async re=>{
            if(re.ok){
                const pendingContacts = await  re.json();
                if(Object.keys(pendingContacts).length ===0)return [];
                for(let contact of pendingContacts.contacts){
                    let contactInfo = await this.getContactInfo(contact.user_id);
                    contact.avatar = contactInfo.avatar;
                }
                return pendingContacts.contacts;
            }
            return []
        }).catch(error => {return [];})
        return response
    }

    async openPedingContactForm(){
        function  acceptContact(contact,token){
            //Tạm
            const response = fetch('/api/v1/contact/accept',{
                method:'POST',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body: JSON.stringify({
                    'contact_user_id':contact.user_id
                })
            }).then(async re=>{
                if(re.ok){
                    return 'success';
                }
                return "error";
            }).catch(error => {return 'error'})
            return response
        }
        async function denyContact(contact,token){
            const data = {
                contact_user_id: contact.user_id,
            }
            return fetch('/api/v1/contact',{
                method:'DELETE',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body:JSON.stringify(data)
            }).then(response=>{
                if(response.ok){
                    return 'success';
                }
                return 'error';
            }).catch(
                ()=> {
                    return 'error';
                }
            )
        }
        function createPendingContactElement(contact){
            console.log('pending contact',contact);
            let image_path = 'assets/images/users/user-dummy-img.jpg';
            if(contact.avatar === null){
                contact.avatar = image_path;
            }
            const li = document.createElement('li');
            li.style.padding = '8px 24px';
            li.style.marginTop= '8px';
            
            li.innerHTML =`
                    <div style="display: flex; justify-content: space-between;gap: 10px; align-items: center;">
                        <div style="display: flex; align-items: center; gap:10px">
                            <img src=${contact.avatar} alt="" style="width: 2.5rem;height: 2.5rem; border-radius: 50%;">
                            <div>
                                <h5 class="font-size-14 m-0">${contact.username}</h5>
                                <span>${contact.email}</span>
                            </div>
                        </div>
                        <div> 
                            <button type="button" class="btn btn-soft-primary btn-sm accept">
                                       <i class="bx bx-plus"></i>
                            </button>
                            <button class="deny btn"> 
                                <i class = 'bx bx-trash '></i>
                            </button>
                        </div>
                    </div>`
            return li;
        }
        


        const contacts = await this.getPendingContacts();
        const ul = this.pendingContactForm.querySelector('.pending-contact-list');
        ul.innerHTML = ''
        for(let contact of contacts){
            let li = createPendingContactElement(contact);
            const btnAccept = li.querySelector('.accept');
            btnAccept.onclick =async ()=>{
                let token = this.store.profile.getToken();
                const message =await acceptContact(contact,token);
                if(message === 'success'){
                    ul.removeChild(li);
                    const newContacts = await this.getContacts()
                    this.renderContacts(newContacts);
                }
            }
            const btnDeny = li.querySelector('.deny');
            btnDeny.onclick = async()=>{
                let token = this.store.profile.getToken();
                const message =await denyContact(contact,token);
                if(message === 'success'){
                    ul.removeChild(li);
                }
            }
            ul.appendChild(li);
        }
        const btnClose1 = this.pendingContactForm.querySelector('.btn.btn-link');
        const btnClose2 = this.pendingContactForm.querySelector('.btn-close.btn-close-white');
        btnClose1.onclick = btnClose2.onclick = ()=>{this.closePendingForm()}

        this.pendingContactForm.classList.add('show');
        document.body.appendChild(this.background);
        document.body.classList.add('modal-open');
        document.body.style.overflow = 'hidden';
        this.pendingContactForm.style.display = 'block';
        this.pendingContactForm.classList.add('show');
    }
    closePendingForm(){
        document.body.classList.remove('modal-open');
        document.body.style.overflow = '';
        document.body.removeChild(this.background)
        this.pendingContactForm.style.display = 'none';
        this.pendingContactForm.classList.remove('show');
    }
// =======================================add more contact in group========================
    async getContactInGroup(groupId){
        // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        let token = this.store.profile.getToken();
        const result = fetch(`/api/v1/group/${groupId}/member`,{
            method:'GET',
            headers:{
                'authorization':`Bearer ${token}`
            }
        }).then(async response => {
            if(!response.ok){
                return [];
            }
            const responseJson = await response.json();

            if(Object.keys(responseJson).length == 0)return [];
            return responseJson.Members;
        }).catch(()=>{
            return [];
        })
        return result;
    }
    async openFormAddMoreContactModel(){
        function createContactElement(contact){
            const li = document.createElement('li');
            li.style.padding = '8px 24px';
            li.style.marginTop= '8px';
            li.innerHTML =`
                    <div style="display: flex;gap: 10px; align-items: center;">
                        <input type="checkbox">
                        <div style="display: flex; align-items: center; gap:10px">
                            <img src=${contact.avatar} alt="" style="width: 2.5rem;height: 2.5rem; border-radius: 50%;">
                            <div> 
                                <h5 class="font-size-14 m-0">${contact.username}</h5>
                                <span>${contact.email}</span>
                            </div>
                        </div>
                    </div>`
            return li;
        }
        async function addMembers(groupId,ids,token){
            return fetch(`/api/v1/group/${groupId}/member`,{
                method:'POST',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body:JSON.stringify({
                    'member_ids':ids,
                })
            })
            .then(async response=>{
                if(response.ok){
                    console.log(await response.json())
                    return 'success';
                }
                return 'error';
            }).catch(()=>{return 'error'})
        }
        
        const currentGroupId = this.store.chat.current_conversation_id;
        console.log(currentGroupId)
        const contactInGroup = await this.getContactInGroup(currentGroupId);
        console.log(contactInGroup)
        // const contacts = await this.getUnJoinContactOfGroup(currentGroupId);
        const contacts = await this.getContacts();
        console.log(contacts);
        const ul = this.addMoreContactModal.querySelector('.unjoin-contacts');
        ul.innerHTML = ''
        let userIds = []
        for(let contact of contacts){
            let li = createContactElement(contact);
            let checkbox = li.querySelector('input');
            if(contactInGroup.find(ct => ct.UserID == contact.user_id)!=null){
                checkbox.disabled = true;
                checkbox.checked = true;
            }
            checkbox.onchange = (e)=>{
                if(checkbox.checked){
                    userIds.push(contact.user_id);
                }
                else{
                    userIds = userIds.filter(id => id!=contact.user_id);
                }
            }
            ul.appendChild(li);
        }

        const btnClose1 = this.addMoreContactModal.querySelector('.btn-close.btn-close-white')
        const btnClose2 = this.addMoreContactModal.querySelector('.btn.btn-link');
        btnClose1.onclick = btnClose2.onclick = ()=>{
            this.closeAddMoreContactForm();
        }

        const searchInput = this.addMoreContactModal.querySelector('#searchMoreContactModal');
        const btnSearch = this.addMoreContactModal.querySelector('#contactSearchbtn-addon');
        btnSearch.onclick = async()=>{
            const email = searchInput.value;
            const contactInfo = await this.searchUnjoinedContact(email);
            console.log(contactInfo);
            ul.innerHTML = '';
            if(contactInfo == null){
                return;
            }
            contactInfo.username = contactInfo.Name;
            contactInfo.email  =contactInfo.Email;
            const li = createContactElement(contactInfo);
            let checkbox = li.querySelector('input');
            if(contactInGroup.find(ct=>ct.UserID == contactInfo.ID) !=null){
                checkbox.checked = true;
                checkbox.disabled = true;
            }
            ul.appendChild(li);
        }
        const btnAdd = this.addMoreContactModal.querySelector('.btn.btn-primary');
        btnAdd.onclick= async ()=>{
            if(userIds.length === 0)return;
            let token = this.store.profile.getToken();
            const message = await addMembers(currentGroupId,userIds,token)
            if(message === 'success'){
                toastr.success('Add successfully','Successfully');
                this.store.chat.renderMemberInGroup()
            }
            else {
                toastr.error('Can not add','Fail');
            }
            console.log(userIds);
            // const
        }

        document.body.appendChild(this.background);
        document.body.classList.add('modal-open');
        document.body.style.overflow = 'hidden';


        this.addMoreContactModal.classList.add('show');
        this.addMoreContactModal.style.display = 'block';
        this.addMoreContactModal.classList.add('show');
    }
    closeAddMoreContactForm(){
        document.body.classList.remove('modal-open');
        document.body.style.overflow = '';
        document.body.removeChild(this.background)
        this.addMoreContactModal.style.display = 'none';
        this.addMoreContactModal.classList.remove('show');
    }
    async searchUnjoinedContact(email){
        return fetch(`/api/v1/user/email/${email}`)
        .then(async response => {
            if(!response.ok) return null;
            const contact  = await response.json()

            if(contact.error == null){
                return contact;
            }
            return null;
        })
        .catch(error =>{
            return null;
        })

    }
    async getUnJoinContactOfGroup(groupId){
        const token = this.store.profile.getToken();
        const output = fetch(`/api/v1/contact/not-in-group/${groupId}`,{
            method:'get',
            headers:{
                'authorization':`Bearer ${token}`
            }
        }).then(async response=> {
            if(!response.ok) return [];
            const contacts = await response.json();
            if(Object.keys(contacts).length === 0){
                return [];
            }
            return contacts.contacts;
        })
        return output;
        // return this.getContacts();//tạm
    }
// ============================================================Render contacts================================
    async renderContacts(contacts){
        async function deleteContact(contact,token){
            const data = {
                contact_user_id: contact.user_id,
            }
            // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
            return fetch('/api/v1/contact',{
                method:'DELETE',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body:JSON.stringify(data)
            }).then(response=>{
                if(response.ok){
                    return 'success';
                }
                return 'error';
            }).catch(
                ()=> {
                    return 'error';
                }
            )
        }
        this.contactDiv.innerHTML = ''
        const ul = document.createElement('ul');
        ul.className = 'list-unstyled';
        let image_path = 'assets/image/users/user-dummy-img.jpg';
        for(let contact of contacts){
            const contactInfo = await this.getContactInfo(contact.user_id)
            if(contactInfo.avatar!=null)
                image_path = contactInfo.avatar;//tạm gọi api phải có luôn avatar
            const li = document.createElement('li');
            li.className = 'contact';
            li.innerHTML = `
                <div>
                    <img src = ${image_path}  style = 'width:2.4rem; height:2.4rem; border-radius:50%!important;'>
                    <span class = 'name' style= 'font-size:medium; font-weight:600;' >${contact.username}</span>
                </div>
                <button class = 'delete-contact' style="border:none; background-color: white">
                    <i class="bx bx-dots-vertical-rounded"></i>
                </button>
                <a class = 'btn-delete'>Delete</a>
            `
            const showDeleteBtn = li.querySelector('.delete-contact');
            const btnDelete = li.querySelector('.btn-delete');
            showDeleteBtn.onclick = ()=>{
                btnDelete.style.display = 'block';
                btnDelete.onclick =async ()=>{
                let token = this.store.profile.getToken();
                const result = await deleteContact(contact,token);
                if(result === 'success'){
                    ul.removeChild(li);
                }
              }

              setTimeout(()=>{
                  document.addEventListener('click',function (event){
                      if(!showDeleteBtn.contains(event.target)){
                        btnDelete.style.display = 'none';
                      }
                      else{
                        btnDelete.style.display = 'block';
                      }
                  })
                  btnDelete.addEventListener('click', function (event) {
                      event.stopPropagation(); // Ngăn không cho sự kiện lan ra ngoài
                  });
              },50)

            }
            // li.onclick = ()=>this.selectContact(contact);//chú ý
            ul.appendChild(li);
        }
        this.contactDiv.appendChild(ul);
    }
    selectContact(contact) {
        console.log(`Contact ${contact.name} được chọn.`);
        // this.store.chat.sele
    }
    async getContacts(){
        let token = this.store.profile.getToken();
        // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        return fetch('/api/v1/contacts',{
            method:'GET',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            },
        }).then(async response => {
            if(!response.ok)return [];
            const contact = await response.json()
            if(Object.keys(contact).length === 0){
                return []
            }
            return contact.contacts;
        }).catch(error=>{
            console.log(error);
            return [];
        })
    }


    // ============================================ Find contact===========================================
    findContactByEmail(){
        async function addContact(contact,token){
            const data = {
                contact_user_id: contact.ID,
            }
            // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
            return fetch('/api/v1/contact',{
                method:'POST',
                headers:{
                    'content-type':'application/json',
                    'authorization':`Bearer ${token}`
                },
                body:JSON.stringify(data)
            }).then(response=>{
                if(response.ok){
                    return 'success';
                }
                return 'error';
            }).catch(
                ()=> {
                    return 'error';
                }
            )
        }
        async function checkIsSent(email,token){
            // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
            return fetch('/api/v1/contact/pending-sent',{
                headers:{
                    'authorization':`Bearer ${token}`
                }
            }).then(async response =>{
                if(!response.ok) return 'error';
                let sentContact = await response.json();
                if(Object.keys(sentContact).length == 0) return false;
                sentContact = sentContact.contacts;
                for(let contact of sentContact){
                    if(contact.email==email)
                        return true;
                }
                return false;
            })
        }
        const emailInput = this.addContactModal.querySelector('#addcontactemail-input');
        const email = emailInput.value;
        const token = this.store.profile.getToken();
        fetch(`/api/v1/user/email/${email}`)
            .then(async response => {
                if(!response.ok) return;
                const contact  = await response.json()

                if(contact.error == null){
                    const previosContact = this.findedContactWrapper.querySelector('.result');
                    if(previosContact !=null && this.findedContactWrapper.contains(previosContact)){
                        this.findedContactWrapper.removeChild(previosContact);
                    }
                    const div = this.createFindedContact(contact);
                    const btnAdd = div.querySelector('.btn.btn-add-contact');
                    const is_sent = await checkIsSent(email,token);
                    if(is_sent == true){
                        btnAdd.innerHTML = this.getCheckSvg();
                        btnAdd.onclick = null;
                    }
                    else{
                        btnAdd.onclick = async ()=>{
                           
                            const result = await addContact(contact,token);
                            if(result === 'success'){
                                btnAdd.innerHTML = this.getCheckSvg();
                            }
                        }
                    }
                    this.findedContactWrapper.appendChild(div);
                }
            })
            .catch(error =>{
                console.log(error);
                const message = this.findedContactWrapper.querySelector('.message');
                message.textContent = 'Đã gửi lời mời';
            })
    }
    getCheckSvg(){
        return `<svg width="20px" height="20px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle opacity="0.5" cx="12" cy="12" r="10" stroke="#1C274C" stroke-width="1.5"/>
        <path d="M8.5 12.5L10.5 14.5L15.5 9.5" stroke="#1C274C" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>`
    }
    createFindedContact(contactInfo){
        console.log(contactInfo)
        const contact = document.createElement('div');
        contact.className = 'result';
        contact.innerHTML = `
           <div  style="display: flex; justify-content: space-between; align-items: center;">
            <div style="display: flex; gap:20px;">
                <img src=${contactInfo.avatar} alt="" style="width: 35px; height: 35px; border-radius: 50%;">
                <div>
                    <p style="margin: 0;">${contactInfo.Name}</p>
                    <p style="margin: 0;">${contactInfo.Email}</p>
                </div>
            </div>
            <div>
                    <span class="btn btn-add-contact">
                        <svg fill="#000000" height="20px" width="20px" version="1.1" id="Capa_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
                             viewBox="0 0 78.509 78.509" xml:space="preserve">
                            <g>
                                <path d="M68.305,51.149h-3.032v-3.031v-6h-6h-5.281c-2.281-2.832-4.785-5.04-7.03-6.974c3.829-3.723,6.22-8.918,6.22-14.668
                                    C53.182,9.186,43.996,0,32.706,0S12.23,9.186,12.23,20.476c0,5.75,2.39,10.945,6.219,14.668
                                    C12.318,40.425,4.205,47.729,4.205,64.26v3h33.708v2.218h6h3.033v3.031v6h6h6.326h6v-6v-3.031h3.032h6v-6V57.15v-6L68.305,51.149
                                    L68.305,51.149z M18.23,20.476C18.23,12.494,24.724,6,32.706,6c7.981,0,14.476,6.494,14.476,14.476
                                    c0,7.449-5.656,13.597-12.897,14.386c-0.072,0.007-0.143,0.016-0.215,0.021c-0.347,0.033-0.698,0.046-1.051,0.054
                                    c-0.097,0.002-0.192,0.01-0.289,0.011c-0.153-0.001-0.303-0.012-0.455-0.017c-0.292-0.009-0.584-0.018-0.871-0.044
                                    c-0.108-0.008-0.215-0.021-0.322-0.031C23.862,34.044,18.23,27.908,18.23,20.476z M22.736,39.369
                                    c0.158-0.137,0.315-0.271,0.472-0.406c2.29,0.981,4.736,1.607,7.242,1.858c0.176,0.02,0.355,0.028,0.532,0.044
                                    c0.469,0.036,0.939,0.062,1.411,0.071c0.105,0.001,0.207,0.016,0.312,0.016c0.078,0,0.154-0.011,0.231-0.012
                                    c0.523-0.009,1.045-0.037,1.566-0.079c0.143-0.013,0.287-0.021,0.428-0.036c2.505-0.246,4.965-0.874,7.271-1.862
                                    c0.155,0.135,0.313,0.27,0.472,0.406c1.415,1.217,2.872,2.48,4.272,3.887v4.862v3.031h-3.033h-6v6v4.11h-27.6
                                    C11.187,49.303,17.297,44.047,22.736,39.369z M68.305,63.478h-9.032v9.031h-6.326v-9.031h-9.033V57.15h9.033v-9.031h6.326v9.031
                                    h9.032V63.478z"/>
                            </g>
                        </svg>
                    </span>
            </div>
        </div>
        `
        return contact;
    }
}
class Profile{
    constructor(store){
        this.store = store;
        this.user_id = null;//tạm
        this.email = null;
        this.profileWrapperDiv = document.querySelector('#pills-user');
        this.settingWrapperDiv = document.querySelector('#pills-setting');
    }
    async init(){
        await this.getProfile();
        this.setInfo();
    }
    getToken(){
        const token = localStorage.getItem('token');
        // console.log(token);
        //tạm
        // let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        if(token == null) {
            window.location.href = '/login';
        }
        return token;   
    }
    deleteToken(){
        localStorage.removeItem('token');
    }
    setInfo(){
        const userInfoImage = this.profileWrapperDiv.querySelector('.rounded-circle.avatar-lg.img-thumbnail')
        const settingImage = this.settingWrapperDiv.querySelector('.rounded-circle.avatar-lg.img-thumbnail.user-profile-image')
        const userInfoName = this.profileWrapperDiv.querySelector('.font-size-16.mb-1.text-truncate');
        const userInfoDetailName = this.profileWrapperDiv.querySelector('.detail-username');
        const userInfoDetailEmail = this.profileWrapperDiv.querySelector('.detail-email');

        const settingName= this.settingWrapperDiv.querySelector('.setting-name');
        const settingEmail = this.settingWrapperDiv.querySelector('.setting-email');

        userInfoDetailName.textContent =  userInfoName.textContent = settingName.textContent = this.username;
        userInfoDetailEmail.textContent = settingEmail.textContent = this.email;
        userInfoImage.src = settingImage.src = this.avatar;
    }
    async getProfile(){
        let token = this.getToken();
        return  fetch('/api/v1/user/me',{
            method:'GET',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            }
        }).then(async response=>{
            if(response.ok){
                const userInfo = await response.json()
                this.email = userInfo.Email;
                this.username = userInfo.Name;
                this.avatar = userInfo.avatar;
                this.user_id = userInfo.ID;
            }
            else{
                //yeu cau dang nhap
                window.location.href = '/login';
            }
        })
    }
}
class Setting{

}

class Socket{
    constructor(store){
        this.store = store;
        this.socket();
    }
    
    socket(){
        const socket = new WebSocket(`ws://${window.location.hostname}:1000/api/v1/ws/${this.store.profile.getToken()}`);
        
        window.addEventListener('beforeunload', () => {
            alert('a')
            socket.close();
        });
        socket.onmessage =async (event) => {
            const message = JSON.parse(event.data);
            console.log(message)
            if(message.func === null){
                toastr('Can not use web chat now','Server error');
                return;
            }   
            
            switch(message.func){
                case 'acceptFriend':
                    this.socketAddFriend(message.user_id);
                    break;
                case 'deleteMessage':
                    this.socketDeleteMessage(message.group_id,message.message_id);
                    break;
                case 'createGroup':
                    this.socketCreateGroup(message.group_id,message.group_name,message.owner_id);
                    break;
                case 'addMember':
                    this.socketAddToGroup(message.group_id,message.group_name,message.owner_id);
                    break;
                case 'sendMessage':
                    this.socketReceiveMessage(message.group_id,message.user_id,message.content);
                    break;
                case 'kickMember':
                    this.socketClickMember(message.group_id,message.user_id);
            }

        }

        socket.onopen = function() {
            
            console.log("WebSocket connection established");
        };

        socket.onerror = function(error) {
            console.log("WebSocket error:", error);
        };
        socket.onclose = function(){
            // socket.close();
        }
    }
    async socketAddFriend(userId){
        const contacts = await this.store.contact.getContacts();
        this.store.contact.renderContacts(contacts);
        this.store.chat.addFavoriteUsers(userId);
    }
    async socketDeleteMessage(groupId,messageId){
        this.store.chat.deleteMessageOfOtherConversation(groupId,messageId);
    }
    async socketCreateGroup(groupId,groupName,ownerId){
        if(this.store.profile.user_id == ownerId) return; // người đó tạo nên không cần add chi.
        this.store.chat.addNewGroupToFavoriteGroups(groupId,groupName);
    }
    async socketAddToGroup(groupId,groupName,ownerId){
        if(this.store.profile.user_id == ownerId) return; // người đó tạo nên không cần add chi.
        this.store.chat.addNewGroupToFavoriteGroups(groupId,groupName);
    }
    async socketReceiveMessage(groupId,userId,content){
        if(this.store.profile.user_id == userId) return;
        this.store.chat.receiveMessage(groupId,userId,content,null);
    }
    socketClickMember(groupId,userId){
        this.store.chat.deleteGroupInFavoriteGroups(groupId)
    }
}
async function main(){
    const store = new Store(); // Khởi tạo trạng thái chung
    const chat = new Chat(store);
    const group = new Group(store);
    const contact = new Contact(store);
    const profile = new Profile(store);
    store.setChat(chat);
    store.setContact(contact)
    store.setGroup(group)
    store.setProfile(profile);
    await profile.init();
    await chat.init();
    await contact.init();
    
    const socket = new Socket(store);

    toastr.options = {
        "closeButton": false, // Không hiện nút đóng
        "progressBar": true,  // Thanh tiến trình hiển thị thời gian tắt
        "positionClass": "toast-top-right", // Vị trí: góc phải trên
        "timeOut": "3000",    // Tự động đóng sau 3 giây
    };

    const btnLogout = document.querySelector('.logout');
    console.log(btnLogout)
    btnLogout.onclick = ()=>{
        profile.deleteToken();
        window.location.href = '/login';
    }
}
main();

// Xóa rồi thêm lại
// thêm thành viên vào nhóm => hiển thị giao diện nhóm mới cho thằng được thêm
// xóa nhóm thành viên trong nhóm => thằng nào bị xóa thì xóa giao diện bên nó. 
// tạo nhóm => thằng nào được thêm lúc tạo thì tạo giao diện nhóm mới cho nó
// kết bạn. gửi tin nhắn đã chấp nhận cho người còn lại.
// Được xóa còn lại bao nhiêu thành viên. Còn 1 thì khi load lại bị mất. mà còn 2 thì nhóm lại thành users.
// vậy thì giờ nó rời nhóm hết thì nhóm sẽ xử lý sao.

// ===================== thông tin gửi lên khi được thêm vào nhóm mới. =================

// tên nhóm
// id nhóm
// tin nhắn cuối cùng
// type: nhóm hay người dùng

// ================================== khi vừa đồng ý kết bạn=================================


//==============================Note mới =================================//
// đã kết bạn rồi thì không gửi lời mời lại được nữa.
// lúc thêm thành viên thì thêm tên nhóm nữa khỏi reload lại.
// load tin nhắn phải thêm tên của người nhắn nữa.
// rời nhóm. chuyển role. admin thì ko cho rời.