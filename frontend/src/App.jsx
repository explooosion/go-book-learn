import React, { useState, useEffect } from "react";
import Login from "./Login";
import ProductManager from "./ProductManager";
import "./App.css";

function App() {
  const [user, setUser] = useState(null);
  const [showLogin, setShowLogin] = useState(false);

  // 當元件載入時，從 localStorage 嘗試讀取 token
  useEffect(() => {
    const token = localStorage.getItem("token");
    const username = localStorage.getItem("username");
    if (token && username) {
      setUser({ token, username });
    }
  }, []);

  const handleLoginSuccess = (loginData) => {
    localStorage.setItem("token", loginData.token);
    localStorage.setItem("username", loginData.username);
    setUser(loginData);
    setShowLogin(false);
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    setUser(null);
  };

  return (
    <div>
      <header className="p-4 bg-gray-800 text-white flex justify-between items-center">
        <h1 className="text-xl font-bold">產品管理系統</h1>
        {user ? (
          <div>
            <span>歡迎，{user.username}</span>
            <button
              onClick={handleLogout}
              className="bg-red-500 px-3 py-1 rounded ml-4"
            >
              登出
            </button>
          </div>
        ) : (
          <button
            onClick={() => setShowLogin(true)}
            className="bg-blue-500 px-3 py-1 rounded"
          >
            登入
          </button>
        )}
      </header>
      {showLogin && !user ? (
        <Login
          onLoginSuccess={handleLoginSuccess}
          onCancel={() => setShowLogin(false)}
        />
      ) : (
        <ProductManager user={user} />
      )}
    </div>
  );
}

export default App;
