<!DOCTYPE html>
<html lang="ar" dir="rtl">
<head>
  <meta charset="UTF-8" />
  <title>بارت جديد - Toxic Rose</title>
  <script src="https://www.gstatic.com/firebasejs/9.22.2/firebase-app-compat.js"></script>
  <script src="https://www.gstatic.com/firebasejs/9.22.2/firebase-auth-compat.js"></script>
  <script src="https://www.gstatic.com/firebasejs/9.22.2/firebase-firestore-compat.js"></script>
  <style>
    body { background: #121212; color: #eee; font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif; padding: 20px; }
    .paragraph { margin-bottom: 30px; }
    .comments-section { margin-top: 10px; border: 1px solid #555; padding: 10px; display: none; }
    .comment { background: #222; color: #eee; margin: 5px 0; padding: 5px; border-radius: 4px; display: flex; justify-content: space-between; }
    .comment button { background: transparent; border: none; color: #f55; cursor: pointer; }
    textarea { width: 100%; height: 60px; }
    button { margin-top: 5px; cursor: pointer; }
    #authButtons button { margin-right: 10px; }
  </style>
</head>
<body>

  <h2>عنوان البارت هنا</h2>

  <div id="authButtons">
    <button id="loginBtn">تسجيل الدخول بـ Google</button>
    <button id="logoutBtn" style="display:none;">تسجيل الخروج</button>
    <span id="userInfo" style="margin-left: 15px;"></span>
  </div>

  <!-- 🔽 ابدأي نسخ الفقرات من هون وعدليها فقط -->
  <div class="paragraph" id="para-1">
    <p>نص الفقرة الأولى من البارت...</p>
    <button onclick="toggleComments('para-1')">تعليقات</button>
    <div class="comments-section" id="comments-para-1">
      <div class="comments-list"></div>
      <textarea placeholder="اكتب تعليقك هنا..." disabled></textarea>
      <button onclick="submitComment('para-1')" disabled>أرسل</button>
    </div>
  </div>

  <div class="paragraph" id="para-2">
    <p>نص الفقرة الثانية من البارت...</p>
    <button onclick="toggleComments('para-2')">تعليقات</button>
    <div class="comments-section" id="comments-para-2">
      <div class="comments-list"></div>
      <textarea placeholder="اكتب تعليقك هنا..." disabled></textarea>
      <button onclick="submitComment('para-2')" disabled>أرسل</button>
    </div>
  </div>
  <!-- 🔼 توقفي هون عند آخر فقرة -->

<script>
  const firebaseConfig = {
    apiKey: "AIzaSyBtTc7yWNfNkG0oVSbpq0V9A6DHTgZoGBM",
    authDomain: "works-rawan.firebaseapp.com",
    projectId: "works-rawan",
    storageBucket: "works-rawan.firebasestorage.app",
    messagingSenderId: "986254083746",
    appId: "1:986254083746:web:17f7db0389c94473f0b9fb"
  };
  firebase.initializeApp(firebaseConfig);

  const auth = firebase.auth();
  const db = firebase.firestore();

  let currentUser = null;

  auth.onAuthStateChanged(user => {
    currentUser = user;
    if(user) {
      loginBtn.style.display = 'none';
      logoutBtn.style.display = 'inline-block';
      userInfo.textContent = `مرحباً، ${user.displayName}`;
      enableCommentBoxes(true);
    } else {
      loginBtn.style.display = 'inline-block';
      logoutBtn.style.display = 'none';
      userInfo.textContent = '';
      enableCommentBoxes(false);
    }
  });

  loginBtn.onclick = () => {
    const provider = new firebase.auth.GoogleAuthProvider();
    auth.signInWithPopup(provider).catch(console.error);
  };
  logoutBtn.onclick = () => { auth.signOut(); };

  function enableCommentBoxes(enable) {
    const textareas = document.querySelectorAll('.comments-section textarea');
    const buttons = document.querySelectorAll('.comments-section button');
    textareas.forEach(ta => ta.disabled = !enable);
    buttons.forEach(btn => btn.disabled = !enable);
  }

  function toggleComments(paraId) {
    const section = document.getElementById('comments-' + paraId);
    if (section.style.display === 'none') {
      section.style.display = 'block';
      loadComments(paraId);
    } else {
      section.style.display = 'none';
    }
  }

  function loadComments(paraId) {
    const commentsList = document.querySelector(`#comments-${paraId} .comments-list`);
    commentsList.innerHTML = 'جاري التحميل...';

    db.collection("comments")
      .where("paragraphId", "==", paraId)
      .orderBy("timestamp", "asc")
      .get()
      .then(querySnapshot => {
        commentsList.innerHTML = '';
        querySnapshot.forEach(doc => {
          const data = doc.data();
          const div = document.createElement('div');
          div.className = 'comment';

          const spanText = document.createElement('span');
          spanText.textContent = data.text;
          div.appendChild(spanText);

          if(currentUser && data.userId === currentUser.uid) {
            const delBtn = document.createElement('button');
            delBtn.textContent = 'حذف';
            delBtn.onclick = () => {
              if(confirm('هل تريد حذف هذا التعليق؟')) {
                db.collection('comments').doc(doc.id).delete();
                loadComments(paraId);
              }
            };
            div.appendChild(delBtn);

            const editBtn = document.createElement('button');
            editBtn.textContent = 'تعديل';
            editBtn.style.marginLeft = '10px';
            editBtn.onclick = () => {
              const newText = prompt('عدّل تعليقك:', data.text);
              if(newText !== null && newText.trim() !== '') {
                db.collection('comments').doc(doc.id).update({ text: newText.trim() });
                loadComments(paraId);
              }
            };
            div.appendChild(editBtn);
          }

          commentsList.appendChild(div);
        });
      })
      .catch(error => {
        commentsList.innerHTML = 'خطأ بتحميل التعليقات';
      });
  }

  function submitComment(paraId) {
    const textarea = document.querySelector(`#comments-${paraId} textarea`);
    const text = textarea.value.trim();
    if (text === '') {
      alert('اكتب تعليق قبل الإرسال');
      return;
    }
    if(!currentUser) {
      alert('لازم تسجل الدخول لتكتب تعليق');
      return;
    }
    db.collection("comments").add({
      paragraphId: paraId,
      text: text,
      userId: currentUser.uid,
      timestamp: firebase.firestore.FieldValue.serverTimestamp()
    })
    .then(() => {
      textarea.value = '';
      loadComments(paraId);
    })
    .catch(() => {
      alert('خطأ بإنشاء التعليق');
    });
  }
</script>

</body>
</html>
