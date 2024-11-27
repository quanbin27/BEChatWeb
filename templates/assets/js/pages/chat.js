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
        this.emailInput = this.addContactModal.querySelector('.find-user');
        this.btnSearchContact = this.addContactModal.querySelector('.modal-footer .btn.btn-primary');
        // this.btnSearchContact.onclick = ()=>alert('Searching')


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
    }

    renderContacts(contacts){
        this.contactDiv.innerHTML = ''
        const ul = document.createElement('ul');
        ul.className = 'list-unstyled';
        for(let contact of contacts){
            const li = document.createElement('li');
            li.className = 'contact';
            li.innerHTML = `
                <img src = ${contact.image_path}  style = 'width:2.4rem; height:2.4rem; border-radius:50%!important;'>
                <span class = 'name' style= 'font-size:medium; font-weight:600;' >${contact.name}</span>
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
        let token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkQXQiOjE3MzMyOTg4MzksInVzZXJfaWQiOiIyIn0.srJB58MbZbN76nxOw3QEPq2-xJkw60Grl9dtugo_EOM'
        fetch('/api/v1/contacts',{
            method:'GET',
            headers:{
                'content-type':'application/json',
                'authorization':`Bearer ${token}`
            },
        }).then(async response => {
            if(!response.ok)return;
            const contact = await response.json()
            return contact;
        })
        // const contacts = [
        //     {
        //         id: 1,
        //         email: 'habc',
        //         image_path:'assets/images/users/avatar-1.jpg',
        //         name:'Nguyễn An',
        //     },
        //     {
        //         id: 2,
        //         email: 'habc',
        //         image_path:'assets/images/users/avatar-1.jpg',
        //         name:'Nguyễn thị bảo',
        //     },
        //     {
        //         id: 3,
        //         email: 'habc',
        //         image_path:'https://s120-ava-talk.zadn.vn/4/8/8/f/4/120/e39bbde1f51d8b1bac7f79ce4510bd7d.jpg',
        //         name:'Tấn',
        //     },
        //     {
        //         id: 10,
        //         email: 'habc',
        //         image_path:'https://s120-ava-talk.zadn.vn/4/8/8/f/4/120/e39bbde1f51d8b1bac7f79ce4510bd7d.jpg',
        //         name:'Chưa nhắn lần nào',
        //     },

        // ]
        // return contacts;
    }
    async getUnJoinContactOfGroup(){
        return this.getContacts();//tạm
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
