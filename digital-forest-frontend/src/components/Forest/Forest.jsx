// src/components/Forest/Forest.jsx
import React, { useState, useRef, useEffect } from 'react';
import Plant from '../Plant/Plant';
import PlantInfo from '../PlantInfo/PlantInfo'; // Импортируем новый компонент
import './Forest.css';

function Forest({ plants, loading, error }) {
  const scrollContainerRef = useRef(null);
  const viewportRef = useRef(null); // Ref для родительского контейнера
  const hideTimeoutRef = useRef(null);

  const [activePlant, setActivePlant] = useState(null); // { plant, position }

  // Очищаем таймер при размонтировании компонента
  useEffect(() => {
    return () => {
      if (hideTimeoutRef.current) {
        clearTimeout(hideTimeoutRef.current);
      }
    };
  }, []);

  const scroll = (direction) => {
    if (scrollContainerRef.current) {
      const scrollAmount = 400;
      scrollContainerRef.current.scrollBy({
        left: direction * scrollAmount,
        behavior: 'smooth'
      });
    }
  };

  const handleMouseEnter = (plant, element) => {
    if (hideTimeoutRef.current) {
      clearTimeout(hideTimeoutRef.current);
    }

    const plantRect = element.getBoundingClientRect();
    const viewportRect = viewportRef.current.getBoundingClientRect();

    // Вычисляем позицию для инфо-блока
    const position = {
      // Центрируем относительно растения
      left: plantRect.left + plantRect.width / 2 - viewportRect.left,
      // Позиционируем ниже уровня земли
      top: plantRect.bottom - viewportRect.top + 30 // 30px ниже основания растения
    };

    setActivePlant({ plant, position });
  };

  const handleMouseLeave = () => {
    hideTimeoutRef.current = setTimeout(() => {
      setActivePlant(null);
    }, 250); // Уменьшаем время для более отзывчивого интерфейса
  };

  const handlePlantInfoMouseEnter = () => {
    // Отменяем таймер, если пользователь навел на PlantInfo
    if (hideTimeoutRef.current) {
      clearTimeout(hideTimeoutRef.current);
    }
  };

  const handlePlantInfoMouseLeave = () => {
    // Запускаем таймер заново при уходе с PlantInfo
    hideTimeoutRef.current = setTimeout(() => {
      setActivePlant(null);
    }, 250);
  };

  const handleWheel = (e) => {
    // Предотвращаем стандартную прокрутку страницы
    e.preventDefault();
    
    if (scrollContainerRef.current) {
      const scrollAmount = e.deltaY > 0 ? 100 : -100; // Прокрутка вниз = вправо, вверх = влево
      scrollContainerRef.current.scrollBy({
        left: scrollAmount
      });
    }
  };


  if (loading) {
    return <div className="forest-status">Загрузка леса...</div>;
  }

  if (error) {
    return <div className="forest-status error">{error}</div>;
  }

  return (
    <div className="forest-container">
      <div className="forest-viewport" ref={viewportRef} onWheel={handleWheel}>
        <div className="forest-scroll-container" ref={scrollContainerRef}>
          <div className="forest-content">
            {plants && plants.length > 0 ? (
              plants.map(plant => (
                <Plant
                  key={plant.id}
                  plant={plant}
                  onMouseEnter={handleMouseEnter}
                  onMouseLeave={handleMouseLeave}
                />
              ))
            ) : (
              <div className="forest-status empty">В лесу пока нет растений. Будь первым!</div>
            )}
          </div>
        </div>

        {/* Рендерим PlantInfo здесь, вне прокручиваемого контейнера */}
        <PlantInfo
          plant={activePlant?.plant}
          position={activePlant?.position}
          onMouseEnter={handlePlantInfoMouseEnter}
          onMouseLeave={handlePlantInfoMouseLeave}
        />
      </div>
      
      {/* Кнопки навигации теперь под лесом */}
      <div className="forest-navigation">
        <button className="scroll-arrow left" onClick={() => scroll(-1)}>◀</button>
        <button className="scroll-arrow right" onClick={() => scroll(1)}>▶</button>
      </div>
    </div>
  );
}

export default Forest;