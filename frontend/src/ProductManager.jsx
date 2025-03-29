import React, { useState, useEffect } from "react";

function ProductManager() {
  const [products, setProducts] = useState([]);
  const [newProduct, setNewProduct] = useState({ name: "", price: "" });
  const [editingProduct, setEditingProduct] = useState(null);

  // 取得所有產品
  const fetchProducts = () => {
    fetch("/api/products")
      .then((res) => res.json())
      .then((data) => setProducts(data))
      .catch((err) => console.error("Error fetching products:", err));
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  // 新增產品
  const handleCreateProduct = () => {
    fetch("/api/products", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: newProduct.name,
        price: parseFloat(newProduct.price),
      }),
    })
      .then((res) => res.json())
      .then(() => {
        fetchProducts();
        setNewProduct({ name: "", price: "" });
      })
      .catch((err) => console.error("Error creating product:", err));
  };

  // 刪除產品
  const handleDeleteProduct = (id) => {
    fetch(`/api/products/${id}`, { method: "DELETE" })
      .then((res) => res.json())
      .then(() => fetchProducts())
      .catch((err) => console.error("Error deleting product:", err));
  };

  // 編輯產品：選擇產品進入編輯模式
  const handleEditProduct = (product) => setEditingProduct(product);

  // 更新產品
  const handleUpdateProduct = () => {
    fetch(`/api/products/${editingProduct.id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: editingProduct.name,
        price: parseFloat(editingProduct.price),
      }),
    })
      .then((res) => res.json())
      .then(() => {
        fetchProducts();
        setEditingProduct(null);
      })
      .catch((err) => console.error("Error updating product:", err));
  };

  return (
    <div className="p-6 max-w-3xl mx-auto">
      <h1 className="text-3xl font-bold mb-4 text-center">產品管理</h1>

      <div className="bg-white p-4 rounded shadow mb-6">
        <h2 className="text-xl font-semibold mb-2">新增產品</h2>
        <div className="flex space-x-2">
          <input
            type="text"
            placeholder="名稱"
            value={newProduct.name}
            onChange={(e) =>
              setNewProduct({ ...newProduct, name: e.target.value })
            }
            className="border p-2 rounded w-full"
          />
          <input
            type="number"
            placeholder="價格"
            value={newProduct.price}
            onChange={(e) =>
              setNewProduct({ ...newProduct, price: e.target.value })
            }
            className="border p-2 rounded w-full"
          />
          <button
            onClick={handleCreateProduct}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
          >
            新增
          </button>
        </div>
      </div>

      <div className="bg-white p-4 rounded shadow">
        <h2 className="text-xl font-semibold mb-2">產品列表</h2>
        {products.length === 0 ? (
          <p>目前無產品資料。</p>
        ) : (
          <ul>
            {products.map((product) => (
              <li
                key={product.id}
                className="flex justify-between items-center mb-2"
              >
                {editingProduct && editingProduct.id === product.id ? (
                  <div className="flex space-x-2 w-full">
                    <input
                      type="text"
                      value={editingProduct.name}
                      onChange={(e) =>
                        setEditingProduct({
                          ...editingProduct,
                          name: e.target.value,
                        })
                      }
                      className="border p-2 rounded w-full"
                    />
                    <input
                      type="number"
                      value={editingProduct.price}
                      onChange={(e) =>
                        setEditingProduct({
                          ...editingProduct,
                          price: e.target.value,
                        })
                      }
                      className="border p-2 rounded w-full"
                    />
                    <button
                      onClick={handleUpdateProduct}
                      className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
                    >
                      更新
                    </button>
                    <button
                      onClick={() => setEditingProduct(null)}
                      className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600"
                    >
                      取消
                    </button>
                  </div>
                ) : (
                  <div className="flex justify-between items-center w-full">
                    <span className="flex-1">
                      {product.name} - ${product.price}
                    </span>
                    <div className="space-x-2">
                      <button
                        onClick={() => handleEditProduct(product)}
                        className="bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600"
                      >
                        編輯
                      </button>
                      <button
                        onClick={() => handleDeleteProduct(product.id)}
                        className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
                      >
                        刪除
                      </button>
                    </div>
                  </div>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

export default ProductManager;
