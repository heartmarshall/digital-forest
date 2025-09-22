// src/utils/imageConverter.js

const PIXEL_SCALE = 10; // Каждый "пиксель" на холсте будет 10x10 реальных пикселей

/**
 * Конвертирует массив цветов (состояние грида) в base64 PNG.
 * @param {string[]} grid - Одномерный массив цветов.
 * @param {number} gridSize - Размер стороны сетки (например, 16 для 16x16).
 * @returns {Promise<string>} - Промис, который разрешается строкой base64 без префикса.
 */
export const gridToBase64 = (grid, gridSize) => {
  return new Promise((resolve, reject) => {
    // Создаем временный canvas, он не будет виден на странице
    const canvas = document.createElement('canvas');
    canvas.width = gridSize * PIXEL_SCALE;
    canvas.height = gridSize * PIXEL_SCALE;
    const ctx = canvas.getContext('2d');

    if (!ctx) {
      return reject(new Error('Не удалось получить 2D контекст canvas'));
    }

    // Проходим по каждому пикселю в нашем состоянии
    grid.forEach((color, index) => {
      const x = (index % gridSize) * PIXEL_SCALE;
      const y = Math.floor(index / gridSize) * PIXEL_SCALE;

      ctx.fillStyle = color;
      ctx.fillRect(x, y, PIXEL_SCALE, PIXEL_SCALE);
    });

    // Экспортируем canvas в data URL (это и есть base64 с префиксом)
    const dataUrl = canvas.toDataURL('image/png');
    // Ваш бэкенд ожидает только саму строку base64, без префикса
    const base64String = dataUrl.replace(/^data:image\/png;base64,/, '');

    resolve(base64String);
  });
};