class Store {
    constructor() {
        this.state = {
            selectedContact: null,
        };
        this.listeners = [];
    }

    // Đăng ký lắng nghe thay đổi
    subscribe(listener) {
        this.listeners.push(listener);
    }

    // Cập nhật trạng thái và thông báo cho listener
    setState(newState) {
        this.state = { ...this.state, ...newState };
        this.listeners.forEach(listener => listener(this.state));
    }

    // Trả về trạng thái hiện tại
    getState() {
        return this.state;
    }
}
const store = new Store(); // Khởi tạo trạng thái chung

class Chat{
    constructor(){
        this.userConversation = document.querySelector('#users-conversation');
        this.btnSend = document.querySelector('.links-list-item button[type="submit"]');
        this.chatinput = document.querySelector('.chat-input');
        this.favouriteUsers= document.querySelector('#favourite-users')
        this.avatar = document.querySelector('.avatar-sm');
        this.profileDetail = document.querySelector('.user-profile-sidebar');
        this.profileAvatar = this.profileDetail.querySelector('.profile-img');
        this.userName = document.querySelector('.user-profile-show')
        this.current_conversation = [];
        this.current_conversation_id = null;
        store.subscribe((state)=>this.renderConversationOfContact(state.selectedContact));
    }

    async init(){
        this.favouriteUserData = await this.getFavouriteUsers();
        if(this.favouriteUserData.length>0){
            this.avatar.src = this.favouriteUserData[0].image_path;
            this.userName.textContent = this.favouriteUserData[0].name;
            this.profileAvatar.src = this.favouriteUserData[0].image_path;
        }
        else{
            //ẩn hết đi
        }
        this.renderFavouriteUsers(this.favouriteUserData);
        const currentConversationId = this.favouriteUserData[0].id;

        this.current_conversation = await this.getCurrentConversationById(currentConversationId);
        this.renderConversation(this.current_conversation);

        this.chatinput.addEventListener('keydown',(event)=>{
            if(this.chatinput.value == '')return;
            if(event.key === 'Enter'){
                this.btnSend.onclick();
            }
        })
        this.btnSend.onclick = (event)=>{
            event.preventDefault();
            if(this.chatinput.value == '')
                return;
            this.current_conversation.push({
                id:null,
                user:'me',
                message:this.chatinput.value
            })
            this.chatinput.value = '';
            this.renderConversation(this.current_conversation);
        }

    }
    sendMessage(messageData){
        // messageData = {
        //     conversationId:1 ,
        //     message: 'xin chào',
        //     // xuống đó tự biết người gửi mà xác định;
        // }
        // fetch()
    }
    async getCurrentConversationById(conversationId){
        // console.log(conversationId)
        //tạm
        const conversation = [
            {
                id:1,//id user
                user:'me',//tên user
                message: 'xin chào bạn nhé', //message
                time: '20:08pm'//giờ
            },
            {
                id:2,
                user:'Nguyễn an',
                message: 'Chào bạn nha',
                time: '20:08pm'//giờ
            },
            {
                id:3,
                user:'me',
                message: 'Một ngày vui vẻ nhé',
                time: '20:08pm'//giờ
            }
        ]
        return conversation;
    }
    async getFavouriteUsers(){
        const users = [
            {
                id:1,//id nhóm
                userId: 1,
                name:'Nguyễn an',// tên nhóm hoặc người dùng
                email: 'testabc@gmail.com', //bỏ
                image_path:'assets/images/users/avatar-1.jpg',
                message: 'hi chào bạn',
                type:'user',
            },
            {
                id:2,
                userId: 2,
                name:'Nguyễn thị bảo',
                email: 'testabc@gmail.com',
                image_path:'assets/images/users/avatar-1.jpg',
                message:'khỏe không',
                type:'user',
            },
            {
                id:3,
                userId: 3,
                name:'Anh tấn sha đô',
                email: 'testabc@gmail.com',
                image_path:'https://s120-ava-talk.zadn.vn/4/8/8/f/4/120/e39bbde1f51d8b1bac7f79ce4510bd7d.jpg',
                message: 'buồn quá đi chơi không mầy',
                type:'user',
            },
            {
                id:4,
                userId: 4,
                name:'Anh quân sha đô',
                email: 'testabc@gmail.com',
                image_path:'https://s120-ava-talk.zadn.vn/0/3/c/6/6/120/7525fac6e9155a92876e5c66245408b8.jpg',
                message:'Cô giao bài tập kìa',
                type:'user',
            },
        ]
        return users;
    }
    ChatComponent(conversation){
        const chatListRight = document.createElement('li');
        if(conversation.user=='me')
            chatListRight.className = 'chat-list right'
        else
            chatListRight.className = 'chat-list left';
        const now = new Date();
        const hours = now.getHours(); // Lấy giờ
        const minutes = now.getMinutes(); // Lấy phút
        const seconds = now.getSeconds(); // Lấy giây
        chatListRight.innerHTML =`
        <div class="conversation-list">
            <div class="user-chat-content">
                <div class="ctext-wrap">
                    <div class="ctext-wrap-content">
                        <p class="mb-0 ctext-content">${conversation.message}</p>
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
        //     <a class="dropdown-item d-flex align-items-center justify-content-between reply-message" href="#" data-bs-toggle="collapse" data-bs-target=".replyCollapse">
        //     Reply <i class="bx bx-share ms-2 text-muted"></i>
        // </a>
        // <a class="dropdown-item d-flex align-items-center justify-content-between" href="#" data-bs-toggle="modal" data-bs-target=".forwardModal">
        //     Forward <i class="bx bx-share-alt ms-2 text-muted"></i>
        // </a>
        // <a class="dropdown-item d-flex align-items-center justify-content-between copy-message" href="#" id="copy-message-1">
        //     Copy <i class="bx bx-copy text-muted ms-2"></i>
        // </a>
        // <a class="dropdown-item d-flex align-items-center justify-content-between" href="#">
        //     Bookmark <i class="bx bx-bookmarks text-muted ms-2"></i>
        // </a>
        // <a class="dropdown-item d-flex align-items-center justify-content-between" href="#">
        //     Mark as Unread <i class="bx bx-message-error text-muted ms-2"></i>
        // </a>
        return chatListRight;
    }
    renderConversationOfContact(state){
        console.log(state)
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
    renderConversation(current_conversation){
        this.userConversation.innerHTML = ''
        for(let messageData of current_conversation){
            const message = this.ChatComponent(messageData);
            this.userConversation.appendChild(message);
        }
    }
    renderFavouriteUsers(favouriteUserData){
        this.favouriteUsers.innerHTML = '';
        for(let userData of favouriteUserData){
            const user = document.createElement('li');
            user.style.cursor = 'pointer';
            user.onmouseover=()=>{
                user.style.backgroundColor = '#d5d5d5'
            }
            user.onmouseleave=()=>{
                user.style.backgroundColor = 'white';
            }
            user.style.marginTop = '10px';
            user.innerHTML = `
                <div style='display:flex; gap:10px; margin-left:10px; align-items:center  '>
                    <img src = ${userData.image_path} style = 'width:2.4rem; height:2.4rem; border-radius:50%!important;'>
                    <div style = 'display:flex; flex-direction:column'>
                        <span style= 'font-size: medium; font-weight:600;' class = 'username'>${userData.name}</span>
                        <span >${userData.message}</span>
                    </div>
                </div>
            `
            user.setAttribute('data-user-id',userData.userId);
            user.setAttribute('data-chat-id',userData.id);
            user.onclick = ()=>{
                this.selectedUser(userData)
            }
            this.favouriteUsers.appendChild(user);

        }
    }
    async selectedUser(userData){
        const user = this.favouriteUsers.querySelector(`li[data-user-id="${userData.userId}"]`);
        this.userName.textContent = userData.name;
        this.profileAvatar.src = userData.image_path;
        this.avatar.src = userData.image_path;
        const currentConversationId = userData.id;
        this.current_conversation = await this.getCurrentConversationById(currentConversationId);
        this.renderConversation(this.current_conversation);
        const currentChat = this.favouriteUsers.querySelector('li.current');
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
}

class Group{
    constructor(){
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
        const friends = [
            {
                id:1,
                name:'Nguyễn an',
                email: 'testabc@gmail.com',
                image_path:'assets/images/users/avatar-1.jpg',
            },
            {
                id:2,
                name:'Nguyễn thị bảo',
                email: 'testabc@gmail.com',
            },
            {
                id:3,
                name:'Anh tấn sha đô',
                email: 'testabc@gmail.com',
                image_path:'https://s120-ava-talk.zadn.vn/4/8/8/f/4/120/e39bbde1f51d8b1bac7f79ce4510bd7d.jpg',
            },
            {
                id:4,
                name:'Anh quân sha đô',
                email: 'testabc@gmail.com',
                image_path:'https://s120-ava-talk.zadn.vn/0/3/c/6/6/120/7525fac6e9155a92876e5c66245408b8.jpg',
            },
        ]
        return friends;
    }
    closeForm(){
        document.body.classList.remove('modal-open');
        document.body.style.overflow = '';
        document.body.removeChild(this.background)
        this.groupModal.style.display = 'none';
        this.groupModal.classList.remove('show');
    }
    openForm(){
        this.background = document.createElement('div');
        this.background.className = 'modal-backdrop fade show';
        document.body.appendChild(this.background);
        document.body.classList.add('modal-open');
        document.body.style.overflow = 'hidden';
        this.groupModal.style.display = 'block';
        this.groupModal.classList.add('show');

        const friendTag = document.querySelector('.card-body.p-2 .simplebar-content .list-unstyled.contact-list');
        console.log(friendTag)
        let currentId = 1;
        const listCheckedFriends = []
        for( let friendData of this.getFriendList()){
            // Tạo phần tử <li>
            const friend = document.createElement('li');

            // Tạo phần tử <div> với class "form-check"
            const formCheckDiv = document.createElement('div');
            formCheckDiv.classList.add('form-check');

            // Tạo phần tử <input> với class "form-check-input"
            const input = document.createElement('input');
            input.type = 'checkbox';

            input.onchange = ()=>{
                if(input.checked){
                    listCheckedFriends.push(friendData.id);
                }
                else
                    listCheckedFriends = listCheckedFriends.filter(fr => fr.id!=friendData.id);

            }

            input.classList.add('form-check-input');
            input.id = `memberCheck${currentId}`;

            // Tạo phần tử <label> với class "form-check-label"
            const label = document.createElement('label');
            label.classList.add('form-check-label');
            label.setAttribute('for', `memberCheck${currentId}`);
            label.textContent = friendData.name; // Gắn tên bạn bè

            // Gắn các phần tử vào nhau
            formCheckDiv.appendChild(input); // Gắn <input> vào <div>
            formCheckDiv.appendChild(label); // Gắn <label> vào <div>
            friend.appendChild(formCheckDiv); // Gắn <div> vào <li>
            friendTag.appendChild(friend); // Gắn <li> vào danh sách <ul>

            // Tăng ID
            currentId++;
        }
        this.btnCreateGroup.onclick = ()=>{
            console.log(listCheckedFriends);
        }
    }

    sendData(){

        //gửi dữ liệu đi;
    }
}

class Contact{
    constructor(){
        this.btnOpenFormAddContact = document.querySelector('.btn.btn-soft-primary.btn-sm');
        this.addContactModal = document.querySelector('#addContact-exampleModal');
        this.emailInput = this.addContactModal.querySelector('#addcontactemail-input');
        this.btnSearchContact = this.addContactModal.querySelector('.modal-footer .btn.btn-primary');


        this.contactDiv = document.querySelector('.sort-contact');
        this.addMoreContactModal = document.querySelector('.modal.fade.addMoreContactModal')
        this.background = document.createElement('div');
        this.background.className = 'modal-backdrop fade show';
        this.btnAddMoreContact = document.querySelector('.add-more-contact');
        this.btnClose1 = this.addMoreContactModal.querySelector('.btn-close.btn-close-white')
        this.btnClose2 = this.addMoreContactModal.querySelector('.btn.btn-link')
    }
    createContactElement(contact){
        const li = document.createElement('li');
        li.style.padding = '8px 24px';
        li.style.marginTop= '8px';
        li.innerHTML =`
                <div style="display: flex;gap: 10px; align-items: center;">
                    <input type="checkbox">
                    <div style="display: flex; align-items: center; gap:10px">
                        <img src=${contact.image_path} alt="" style="width: 2.5rem;height: 2.5rem; border-radius: 50%;">
                        <h5 class="font-size-14 m-0">${contact.name}</h5>
                    </div>
                </div>`
        return li;
    }
    async openFormAddMoreContactModel(){
        const userChat = document.querySelector('.user-chat.w-100.overflow-hidden')
        const currentGroupId = userChat.getAttribute('chat-id');
        const contacts = await this.getUnJoinContactOfGroup(currentGroupId);
        const ul = this.addMoreContactModal.querySelector('.unjoin-contacts');
        ul.innerHTML = ''
        for(let contact of contacts){
            let li = this.createContactElement(contact);
            ul.appendChild(li);
        }
        this.addMoreContactModal.classList.add('show')
        document.body.appendChild(this.background);
        document.body.classList.add('modal-open');
        document.body.style.overflow = 'hidden';
        this.addMoreContactModal.style.display = 'block';
        this.addMoreContactModal.classList.add('show');
    }
    closeForm(){
        document.body.classList.remove('modal-open');
        document.body.style.overflow = '';
        document.body.removeChild(this.background)
        this.addMoreContactModal.style.display = 'none';
        this.addMoreContactModal.classList.remove('show');
    }
    seachUnjoinedContact(e){
        const ul = this.addMoreContactModal.querySelector('.unjoin-contacts');
        const list = ul.querySelectorAll('li');
        for(const li of list){
            const name = li.querySelector('.font-size-14.m-0').textContent.toUpperCase();
            if(name.indexOf(e.target.value.toUpperCase())>-1){
                li.style.display = ''
            }
            else li.style.display = 'none';
        }
    }
    async init(){
        const contacts = await this.getContacts()
        console.log(contacts);
        this.renderContacts(contacts);
        this.btnAddMoreContact.onclick = ()=>{
            this.openFormAddMoreContactModel();
        }
        console.log(this.btnClose1)
        this.btnClose1.onclick = this.btnClose2.onclick = ()=>{
            this.closeForm();
        }
        const search = this.addMoreContactModal.querySelector('#searchMoreContactModal');
        console.log(search)
        search.addEventListener('input',(e)=>{
            this.seachUnjoinedContact(e)
        })

        this.btnSearchContact.onclick = ()=>{
            this.findContactByEmail();
        }

    }

    renderContacts(contacts){
        console.log(contacts)
        this.contactDiv.innerHTML = ''
        const ul = document.createElement('ul');
        ul.className = 'list-unstyled';
        let image_path = 'https://s120-ava-talk.zadn.vn/4/8/8/f/4/120/e39bbde1f51d8b1bac7f79ce4510bd7d.jpg';
        for(let contact of contacts){
            const li = document.createElement('li');
            li.className = 'contact';
            li.innerHTML = `
                <img src = ${image_path}  style = 'width:2.4rem; height:2.4rem; border-radius:50%!important;'>
                <span class = 'name' style= 'font-size:medium; font-weight:600;' >${contact.username}</span>
            `
            li.onclick = ()=>this.selectContact(contact);
            ul.appendChild(li);
        }
        this.contactDiv.appendChild(ul);
    }
    selectContact(contact) {
        console.log(`Contact ${contact.name} được chọn.`);
        store.setState({ selectedContact: contact });
    }
    async getContacts(){
        // alert('a')
        let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        return fetch('/api/v1/contacts',{
            method:'GET',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            },
        }).then(async response => {
            if(!response.ok)return;
            const contact = await response.json()
            // console.log(contact.contacts);
            return contact.contacts;
        })
    }
    async getUnJoinContactOfGroup(){
        return this.getContacts();//tạm
    }
    findContactByEmail(){
        const email = this.emailInput.value;

        //console.log(email);
    }
    createFindedContact(contactInfo){
        const contact = document.createElement('div');
        contact.innerHTML = `
           <div class="user" style="display: flex; justify-content: space-between; align-items: center;">
            <div style="display: flex; gap:20px;">
                <img src="assets/images/users/avatar-2.jpg" alt="" style="width: 35px; height: 35px; border-radius: 50%;">
                <div>
                    <p style="margin: 0;">${contactInfo.username}</p>
                    <p style="margin: 0;">${contactInfo.email}</p>
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

}
class Setting{

}
async function main(){
    const chat = new Chat();
    const group = new Group();
    const contact = new Contact();
    await chat.init();
    await contact.init();
}
main();
