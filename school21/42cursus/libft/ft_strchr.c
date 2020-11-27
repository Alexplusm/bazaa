/* ************************************************************************** */
/*                                                                            */
/*                                                        :::      ::::::::   */
/*   ft_strchr.c                                        :+:      :+:    :+:   */
/*                                                    +:+ +:+         +:+     */
/*   By: cdeon <marvin@42.fr>                       +#+  +:+       +#+        */
/*                                                +#+#+#+#+#+   +#+           */
/*   Created: 2020/11/27 09:13:53 by cdeon             #+#    #+#             */
/*   Updated: 2020/11/27 09:13:56 by cdeon            ###   ########.fr       */
/*                                                                            */
/* ************************************************************************** */

char *ft_strchr(const char *s, int c)
{
    char *ptr;

    ptr = (char *)s;
    while(*ptr != '\0')
    {
        if (*ptr == c)
            return ptr;
        ptr++;
    }
    return (c == '\0') ? ptr : 0;
}
