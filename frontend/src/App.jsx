import "./App.css";
import React, { useState } from "react";

import Login from "./Login";
import ProductManager from "./ProductManager"; // 這是你的產品管理頁面

function App() {
  const [user, setUser] = useState(null);

  const handleLoginSuccess = (loginData) => {
    // 這裡你可以存取 token 或其他登入資訊
    setUser(loginData);
  };

  return (
    <div>
      {user ? (
        <ProductManager user={user} />
      ) : (
        <Login onLoginSuccess={handleLoginSuccess} />
      )}
    </div>
  );
}

export default App;
