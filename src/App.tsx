import React, { useEffect, useState } from "react";
import db from "./firebaseFirestore";
import { doc, getDoc } from "firebase/firestore";

const App = () => {
  const [partContent, setPartContent] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function fetchPart() {
      try {
        const docRef = doc(db, "parts", "part1"); // ← غيّري المعرف حسب الحاجة
        const docSnap = await getDoc(docRef);

        if (docSnap.exists()) {
          const data = docSnap.data();
          setPartContent(data?.text || "⚠️ المحتوى فارغ.");
        } else {
          setPartContent("❌ لا يوجد محتوى.");
        }
      } catch (error) {
        console.error("خطأ أثناء تحميل البارت:", error);
        setPartContent("⚠️ فشل في تحميل البارت.");
      } finally {
        setLoading(false);
      }
    }

    fetchPart();
  }, []);

  // تقسيم الفقرات بدقة أكثر
  const paragraphs = partContent.trim().split(/\n{2,}/);

  return (
    <div
      style={{
        padding: "2rem",
        fontFamily: "sans-serif",
        backgroundColor: "#111",
        color: "#eee",
        minHeight: "100vh"
      }}
    >
      {loading ? (
        <div style={{ textAlign: "center", fontSize: "1.2rem" }}>
          <span className="spinner" style={{ display: "inline-block", marginRight: "0.5rem" }}>
            ⏳
          </span>
          جاري تحميل البارت...
        </div>
      ) : (
        paragraphs.map((para, index) => (
          <div
            key={index}
            style={{
              marginBottom: "2rem",
              position: "relative",
              paddingRight: "2.5rem"
            }}
          >
            <p style={{ lineHeight: "1.8", fontSize: "1rem" }}>{para}</p>

            {/* زر تعليق جانبي */}
            <div
              onClick={() => alert(`فتح التعليقات للفقرة رقم ${index + 1}`)} // ← استبدليه لاحقاً بربط Firebase
              style={{
                position: "absolute",
                top: 0,
                right: 0,
                background: "#444",
                borderRadius: "50%",
                width: "1.5rem",
                height: "1.5rem",
                textAlign: "center",
                lineHeight: "1.5rem",
                fontSize: "0.8rem",
                cursor: "pointer",
                transition: "background 0.3s"
              }}
              title="تعليقات"
              onMouseOver={(e) => (e.currentTarget.style.background = "#666")}
              onMouseOut={(e) => (e.currentTarget.style.background = "#444")}
            >
              💬
            </div>
          </div>
        ))
      )}
    </div>
  );
};

export default App;
