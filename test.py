from ucimlrepo import fetch_ucirepo 
import pandas as pd
import numpy as np
import seaborn as sns
import matplotlib.pyplot as plt

from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import accuracy_score, classification_report, confusion_matrix

# 获取 UCI 数据集
bone_marrow_transplant_children = fetch_ucirepo(id=565) 

# 数据集的特征（X）和目标变量（y）
X = bone_marrow_transplant_children.data.features 
y = bone_marrow_transplant_children.data.targets 

# 打印数据集的元信息
print(bone_marrow_transplant_children.metadata) 

# 打印变量信息（数据的列名、类型等）
print(bone_marrow_transplant_children.variables) 

# 查看数据基本信息
print(X.info())  # 查看是否有缺失值
print(X.describe())  # 统计描述

# 确保 y 是 Series
y = y.squeeze()

# 处理目标变量（如果 y 是分类变量）
if y.dtype == 'object':
    y = y.map({'Yes': 1, 'No': 0})  # 假设 "Yes" 是 1，"No" 是 0

print(y.value_counts())  # 查看目标变量分布

# 获取所有数值列
numeric_columns = X.select_dtypes(include=[np.number]).columns

# 只对数值列填充缺失值（使用 .loc 避免 SettingWithCopyWarning）
X.loc[:, numeric_columns] = X[numeric_columns].fillna(X[numeric_columns].median())

# 处理类别变量（独热编码 One-Hot Encoding）
X = pd.get_dummies(X, drop_first=True)  # 这会把分类变量转换为数值

# 标准化数据
scaler = StandardScaler()
X_scaled = scaler.fit_transform(X)

# 80% 训练集, 20% 测试集
X_train, X_test, y_train, y_test = train_test_split(X_scaled, y, test_size=0.2, random_state=42)

# 创建逻辑回归模型
model = LogisticRegression()

# 训练模型
model.fit(X_train, y_train)

# 预测测试集
y_pred = model.predict(X_test)

# 计算准确率
accuracy = accuracy_score(y_test, y_pred)
print(f"模型准确率: {accuracy:.4f}")

# 输出分类报告
print("分类报告：")
print(classification_report(y_test, y_pred))

# 计算混淆矩阵
cm = confusion_matrix(y_test, y_pred)

# 绘制热力图
plt.figure(figsize=(5, 4))
sns.heatmap(cm, annot=True, fmt="d", cmap="Blues")
plt.xlabel("Predicted Label")
plt.ylabel("True Label")
plt.title("Confusion Matrix")
plt.show()
