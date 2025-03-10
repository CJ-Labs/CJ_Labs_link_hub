import { greet } from './index';

test('greet function', () => {
    console.log("djashdjadaj")
    expect(greet('World')).toBe('Hello, World!');
}); 