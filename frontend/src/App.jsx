import React, { useState, useEffect } from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import Login from "./Login";
import ProductManager from "./ProductManager";
import "./App.css";

function App() {
  const [user, setUser] = useState(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    const username = localStorage.getItem("username");
    const role = localStorage.getItem("role");
    if (token && username) {
      setUser({ token, username, role });
    }
  }, []);

  const handleLoginSuccess = (loginData) => {
    localStorage.setItem("token", loginData.token);
    localStorage.setItem("username", loginData.username);
    localStorage.setItem("role", loginData.role);
    setUser(loginData);
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    localStorage.removeItem("role");
    setUser(null);
  };

  const handleTokenRefresh = (newToken) => {
    localStorage.setItem("token", newToken);
    setUser((prev) => ({ ...prev, token: newToken }));
  };

  return (
    <Router>
      <header className="p-4 bg-gray-800 text-white flex justify-between items-center">
        <h1 className="text-xl font-bold">產品管理系統</h1>
        {user ? (
          <div>
            <span>
              歡迎，{user.username + (user.role ? ` (${user.role})` : "")}
            </span>
            <button
              onClick={handleLogout}
              className="bg-red-500 px-3 py-1 rounded ml-4"
            >
              登出
            </button>
          </div>
        ) : (
          <button onClick={() => {}} className="bg-blue-500 px-3 py-1 rounded">
            {/* 未登入時，使用者可自行點選連結進入登入頁面 */}
            <a href="/login" className="text-white no-underline">
              登入
            </a>
          </button>
        )}
      </header>
      <Routes>
        <Route
          path="/login"
          element={<Login onLoginSuccess={handleLoginSuccess} />}
        />
        {/* 產品列表頁面永遠公開，但內部操作會依 token 來控制 */}
        <Route
          path="/products"
          element={
            <ProductManager user={user} onTokenRefresh={handleTokenRefresh} />
          }
        />
        {/* 預設導向 /products */}
        <Route path="*" element={<Navigate to="/products" />} />
      </Routes>
    </Router>
  );
}

export default App;
