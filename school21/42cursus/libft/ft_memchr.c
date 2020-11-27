/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   ft_memchr.c                                        :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: cdeon <marvin@42.fr>                       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2020/11/27 09:12:04 by cdeon             #+#    #+#             */
/*   Updated: 2020/11/27 09:12:06 by cdeon            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

void *ft_memchr(const void *s, int c, int n)
{
    char *ptr;

    ptr = (char *)s;
    while (n-- > 0 && *ptr != c)
        ptr++;
    return (n == -1) ? 0 : ptr;
}
