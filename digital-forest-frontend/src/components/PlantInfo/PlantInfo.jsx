// src/components/PlantInfo/PlantInfo.jsx
import './PlantInfo.css';

function PlantInfo({ plant, position, onMouseEnter, onMouseLeave }) {
  if (!plant || !position) {
    return null;
  }

  const formattedDate = new Date(plant.createdAt).toLocaleDateString();

  // Стили для позиционирования всего блока
  const style = {
    left: `${position.left}px`,
    top: `${position.top}px`,
  };

  return (
    <div
      className="plant-info-wrapper"
      style={style}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <div className="dotted-line"></div>
      <div className="info-box">
        <p>Автор: {plant.author}</p>
        <p>Дата: {formattedDate}</p>
      </div>
    </div>
  );
}

export default PlantInfo;