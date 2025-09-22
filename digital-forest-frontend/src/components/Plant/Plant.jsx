// src/components/Plant/Plant.jsx
import './Plant.css';

function Plant({ plant, onMouseEnter, onMouseLeave }) {
  // Передаем в обработчики само растение и DOM-элемент img
  const handleMouseEnter = (event) => {
    onMouseEnter(plant, event.currentTarget);
  };

  return (
    <div className="plant-container">
      <img
        className="plant-image"
        src={`data:image/png;base64,${plant.imageData}`}
        alt={`Пиксель-арт от ${plant.author}`}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={onMouseLeave}
      />
    </div>
  );
}

export default Plant;